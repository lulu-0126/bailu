/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	corev1 "k8s.io/api/core/v1"

	cloudedgev1 "igitlab.lenovo.com/lecp/mec/bailu/api/v1"

	bailuclient "igitlab.lenovo.com/lecp/mec/bailu/client/clientset/versioned"
)

// NodeManagementReconciler reconciles a NodeManagement object
type NodeManagementReconciler struct {
	nmclient  bailuclient.Interface
	clientset *kubernetes.Clientset
	client.Client
	Scheme *runtime.Scheme
}

func (r *NodeManagementReconciler) initDrainer(config *rest.Config) error {
	r.clientset, _ = kubernetes.NewForConfig(config)
	return nil
}

//+kubebuilder:rbac:groups=cloud.edge,resources=nodemanagements,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud.edge,resources=nodemanagements/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud.edge,resources=nodemanagements/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;create

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NodeManagement object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *NodeManagementReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	klog.Info("Reconciling NodeManagement")
	klog.Infof("nodeSyncWorker start")
	defer klog.Infof("nodeSyncWorker end")
	nodelist := corev1.NodeList{}
	if err := r.List(ctx, &nodelist); err != nil {
		klog.Error("unable to list Node while nodeSync", err)
	}
	for _, node := range nodelist.Items {
		m := &cloudedgev1.NodeManagement{}
		namespacedName := types.NamespacedName{
			Name:      node.Name,
			Namespace: node.Namespace,
		}
		if err := r.Get(ctx, namespacedName, m); err != nil {
			if client.IgnoreNotFound(err) == nil {
				klog.Info("NodeManagement not found, creating a new one")
				if err := r.createNodeManagement(ctx, node.Name); err != nil {
					klog.Error("unable to create NodeManagement:", err)
				}
			} else {
				klog.Error("unable to fetch NodeManagement:", err)
			}
		}
	}

	nmlist, err := r.nmclient.CloudV1().NodeManagements().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error("unable to list NodeManagement", err)
	}
	for _, nm := range nmlist.Items {
		node, err := r.fetchNode(ctx, nm.Name)
		if err != nil && client.IgnoreNotFound(err) == nil {
			klog.Error("Found orphan NodeManagement: ", nm.Name)
			if nm.Status.NodeStatus == cloudedgev1.NodeStatusUnknown {

			} else {
				//nm.Status.NodeStatus = cloudedgev1.NodeStatusUnknown
				token := r.generictoken(ctx, r.clientset)
				//nm.Status.Token = token
				err := r.updateMaintenanceStatus(context.TODO(), &nm, cloudedgev1.NodeStatusUnknown, token)
				if err != nil {
					klog.Error("unable to update NodeManagement:", err)
				}
			}
		} else {
			status := getNodeStatus(node)
			nm.Status.NodeStatus = cloudedgev1.NodeStatusType(status)
			//nm.Status.LastUpdate = metav1.Now()
			nm.Status.NodeKind = getNodeRoles(node)
			nm.Status.NodeVersion = getNodeVersion(node)
			nm.Status.NodeIP = getNodeIP(node)
			nm.Status.NodeArc = getNodeArchitecture(node)
			// TODO:add labels to node
			//labels := MergeMaps(node.Labels, nm.Labels)
			//node.Labels = labels
			//_, err = r.clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
			//if err != nil {
			//	klog.Errorf("unable to update Node with [%s] status", status, err)
			//  panic(err.Error())
			//}
			//fmt.Println("Labels 已成功写入节点")
			err := r.updateMaintenanceStatus(context.TODO(), &nm, nm.Status.NodeStatus, nm.Status.Token)
			klog.Infof("update NodeManagement: %s", nm.Status.NodeStatus)
			klog.Infof("update NodeManagement: %s", node.Name)
			if err != nil {
				klog.Error("unable to update222 NodeManagement:", err)
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *NodeManagementReconciler) fetchNode(ctx context.Context, name string) (*corev1.Node, error) {
	node := &corev1.Node{}
	if err := r.Get(ctx, types.NamespacedName{Name: name}, node); err != nil {
		klog.Error("unable to fetch Node:", err)
		return nil, err
	}
	return node, nil
}

func (r *NodeManagementReconciler) updateMaintenanceStatus(ctx context.Context, m *cloudedgev1.NodeManagement, status cloudedgev1.NodeStatusType, token string) error {
	m.Status.NodeStatus = status
	m.Status.Token = token
	m.Status.LastUpdate = metav1.Now()
	if err := r.Status().Update(ctx, m); err != nil {
		klog.Errorf("unable to update NodeManagement with [%s] status", status, err)
		return err
	}
	return nil
}

func (r *NodeManagementReconciler) handlerNodeCreate(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	klog.Infof("handlerNodeCreate: %s", e.Object.GetName())
	node := e.Object.(*corev1.Node)

	m := &cloudedgev1.NodeManagement{}
	namespacedName := types.NamespacedName{
		Name:      node.Name,
		Namespace: node.Namespace,
	}
	if err := r.Get(ctx, namespacedName, m); err != nil {
		if client.IgnoreNotFound(err) == nil {
			klog.Info("NodeManagement not found, creating a new one")
			if err := r.createNodeManagement(ctx, node.Name); err != nil {
				klog.Error("unable to create NodeManagement:", err)
			}
		} else {
			klog.Error("unable to fetch NodeManagement:", err)
		}
	} else {
		klog.Info("NodeManagement exist, skip creating")
	}
}

func (r *NodeManagementReconciler) createNodeManagement(ctx context.Context, name string) error {
	m := &cloudedgev1.NodeManagement{
		TypeMeta: metav1.TypeMeta{
			Kind: "NodeManagement",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	if err := r.Create(ctx, m); err != nil {
		klog.Error("unable to create NodeManagement:", err)
		return err
	}
	return nil
}

var nodeIp string

func (r *NodeManagementReconciler) generictoken(ctx context.Context, clientset *kubernetes.Clientset) string {
	//创建新的 bootstrap token
	token, err := createBootstrapToken()
	if err != nil {
		panic(err)
	}

	// 创建新的 bootstrap token secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("bootstrap-token-%s", token[:6]),
			Namespace: "kube-system",
		},
		Type: "bootstrap.kubernetes.io/token",
		Data: map[string][]byte{
			"token-id":                       []byte(token[:6]),
			"token-secret":                   []byte(token[6:22]),
			"usage-bootstrap-authentication": []byte("true"),
			"usage-bootstrap-signing":        []byte("true"),
			"auth-extra-groups":              []byte("system:bootstrappers:kubeadm:default-node-token"),
		},
	}
	// 使用 Kubernetes API 将新的 bootstrap token secret 添加到 Kubernetes 集群中
	_, err = clientset.CoreV1().Secrets("kube-system").Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	//host := config.Host
	AuthToken := token[:6] + "." + token[6:22]
	fmt.Println("Successfully created bootstrap token")
	nodelist := corev1.NodeList{}
	if err := r.List(ctx, &nodelist); err != nil {
		klog.Error("unable to list Node while nodeSync", err)
	}
	for _, node := range nodelist.Items {
		GetNodeRole := getNodeRolesForArray(&node)
		if StringInArray("master", GetNodeRole) || StringInArray("control-plane", GetNodeRole) {
			nodeIp = getNodeIP(&node)
			break
		}
	}
	//// 获取 CA 证书
	//caCertData, err := getCACertData(ctx, clientset)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// 计算 CA 证书的哈希值
	//caCertHash := calculateCertHash(caCertData)
	//klog.Infof("Discovery Token CA Cert Hash: sha256:%s\n", caCertHash)

	//生成加入集群的命令
	//joinCommand := fmt.Sprintf("kubeadm join %s --token %s --discovery-token-unsafe-skip-ca-verification", apiServerHost, AuthToken)
	joinCommand := fmt.Sprintf("%s:6443 --token %s --discovery-token-unsafe-skip-ca-verification", nodeIp, AuthToken)
	klog.Infof("Join Command: %s\n", joinCommand)
	return joinCommand
}

func createBootstrapToken() (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

// 获取 CA 证书数据
func getCACertData(ctx context.Context, clientset *kubernetes.Clientset) ([]byte, error) {
	// 获取 Secret 对象
	secret, err := clientset.CoreV1().Secrets("kube-public").Get(ctx, "extension-apiserver-authentication", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// 获取 CA 证书数据
	caCertData, found := secret.Data["client-ca-file"]
	if !found {
		return nil, fmt.Errorf("CA certificate data not found")
	}

	return caCertData, nil
}

// 计算证书哈希值
func calculateCertHash(certData []byte) string {
	hash := sha256.Sum256(certData)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (r *NodeManagementReconciler) handlerNodeUpdate(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	klog.Infof("handlerNodeUpdate: %s %s", e.ObjectOld.GetName(), e.ObjectNew.GetName())
	ctx := context.Background()
	nodelist := corev1.NodeList{}
	if err := r.List(ctx, &nodelist); err != nil {
		klog.Error("unable to list Node while nodeSync", err)
	}
	nmlist, err := r.nmclient.CloudV1().NodeManagements().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error("unable to list NodeManagement", err)
	}
	for _, nm := range nmlist.Items {
		node, err := r.fetchNode(ctx, nm.Name)
		if err != nil && client.IgnoreNotFound(err) == nil {
			klog.Error("Found orphan NodeManagement: ", nm.Name)
			if nm.Status.NodeStatus == "unknown" {

			} else {
				nm.Status.NodeStatus = cloudedgev1.NodeStatusUnknown
				//nm.Status.LastUpdate = metav1.Now()
				err := r.updateMaintenanceStatus(context.TODO(), &nm, nm.Status.NodeStatus, nm.Status.Token)
				if err != nil {
					klog.Error("unable to update NodeManagement:", err)
				}
			}
		} else {
			status := getNodeStatus(node)
			nm.Status.NodeStatus = cloudedgev1.NodeStatusType(status)
			nm.Status.LastUpdate = metav1.Now()
			nm.Status.NodeKind = getNodeRoles(node)
			nm.Status.NodeVersion = getNodeVersion(node)
			nm.Status.NodeIP = getNodeIP(node)
			nm.Status.NodeArc = getNodeArchitecture(node)
			klog.Infof("update NodeManagement: %s", nm.Status.NodeStatus)
			klog.Infof("update NodeManagement: %s", node.Name)
			err := r.updateMaintenanceStatus(context.TODO(), &nm, nm.Status.NodeStatus, nm.Status.Token)
			if err != nil {
				klog.Error("unable to update NodeManagement:", err)
			}
		}
	}
}

type NodeType string

//fetch node status
const (
	// NodeReady means kubelet is healthy and ready to accept pods.
	NodeReady    NodeType = "Ready"
	NotNodeReady NodeType = "NotReady"
)

func getNodeStatus(node *corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
			return string(NodeReady)
		}
	}
	return string(NotNodeReady)
}

// 获取节点 IP 地址
func getNodeIP(node *corev1.Node) string {
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			return address.Address
		}
	}
	return ""
}

// 获取节点版本
func getNodeVersion(node *corev1.Node) string {
	return node.Status.NodeInfo.KubeletVersion
}

// 获取节点架构
func getNodeArchitecture(node *corev1.Node) string {
	return node.Status.NodeInfo.Architecture
}

// 获取节点角色
func getNodeRoles(node *corev1.Node) string {
	roles := []string{}
	for key, _ := range node.Labels {
		if key == "node-role.kubernetes.io/master" {
			roles = append(roles, "master")
		}
		if key == "node-role.kubernetes.io/control-plane" {
			roles = append(roles, "control-plane")
		}
		// 如果有其他自定义的节点角色标签，可以在此处添加相应的逻辑
	}
	strSlice := roles
	str := strings.Join(strSlice, ",")
	return str
}

//获取节点角色
func getNodeRolesForArray(node *corev1.Node) []string {
	roles := []string{}
	for key, _ := range node.Labels {
		if key == "node-role.kubernetes.io/master" {
			roles = append(roles, "master")
		}
		if key == "node-role.kubernetes.io/control-plane" {
			roles = append(roles, "control-plane")
		}
		// 如果有其他自定义的节点角色标签，可以在此处添加相应的逻辑
	}
	return roles
}

// StringInArray 判断字符串是否在字符串数组中
func StringInArray(str string, array []string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}

func MergeMaps(map1, map2 map[string]string) map[string]string {
	merged := make(map[string]string)

	// 将 map1 的键值对添加到 merged 中
	for k, v := range map1 {
		merged[k] = v
	}

	// 将 map2 的键值对添加到 merged 中
	for k, v := range map2 {
		merged[k] = v
	}

	return merged
}


func (r *NodeManagementReconciler) handlerNodeDelete(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	klog.Infof("handlerNodeDelete: %s", e.Object.GetName())
	node := e.Object.(*corev1.Node)

	m := &cloudedgev1.NodeManagement{
		TypeMeta: metav1.TypeMeta{
			Kind: "NodeManagement",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: node.Name,
		},
	}
	if err := r.Delete(ctx, m); client.IgnoreNotFound(err) != nil {
		klog.Error("unable to delete NodeManagement:", err)
	}
}

func (r *NodeManagementReconciler) handlerNodeGeneric(e event.GenericEvent, q workqueue.RateLimitingInterface) {
	// TODO: 什么场景下会有GenericEvnet?
	klog.Infof("handlerNodeGeneric: %s", e.Object.GetName())
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeManagementReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.initDrainer(mgr.GetConfig())
	r.nmclient = bailuclient.NewForConfigOrDie(mgr.GetConfig())

	go func() {
		<-mgr.Elected()
		klog.Infof("NodeManagement controller elected as leader, run nodeSyncWorker")
		//r.nodeSyncWorker()
		//go wait.Until(r.nodeSyncWorker, NodeSyncInterval, nil)
	}()

	funcs := handler.Funcs{
		CreateFunc:  r.handlerNodeCreate,
		UpdateFunc:  r.handlerNodeUpdate,
		DeleteFunc:  r.handlerNodeDelete,
		GenericFunc: r.handlerNodeGeneric,
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudedgev1.NodeManagement{}).
		Watches(&source.Kind{Type: &corev1.Node{}}, &funcs).
		Complete(r)
}

// writer implements io.Writer interface as a pass-through for klog.
type writer struct {
	logFunc func(args ...interface{})
}

// Write passes string(p) into writer's logFunc and always returns len(p)
func (w writer) Write(p []byte) (n int, err error) {
	w.logFunc(string(p))
	return len(p), nil
}
