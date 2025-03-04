package config

import (
	"encoding/pem"
	"os"
	"sync"

	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
	"github.com/kubeedge/kubeedge/pkg/client/clientset/versioned"
	syncinformer "github.com/kubeedge/kubeedge/pkg/client/informers/externalversions/reliablesyncs/v1alpha1"
	synclister "github.com/kubeedge/kubeedge/pkg/client/listers/reliablesyncs/v1alpha1"
)

var Config Configure
var once sync.Once

type Configure struct {
	v1alpha1.CloudHub
	KubeAPIConfig *v1alpha1.KubeAPIConfig
	Ca            []byte
	CaKey         []byte
	Cert          []byte
	Key           []byte
}

func InitConfigure(hub *v1alpha1.CloudHub) {
	once.Do(func() {
		if len(hub.AdvertiseAddress) == 0 {
			klog.Exit("AdvertiseAddress must be specified!")
		}

		Config = Configure{
			CloudHub: *hub,
		}

		ca, err := os.ReadFile(hub.TLSCAFile)
		if err == nil {
			block, _ := pem.Decode(ca)
			ca = block.Bytes
			klog.Info("Succeed in loading CA certificate from local directory")
		}

		caKey, err := os.ReadFile(hub.TLSCAKeyFile)
		if err == nil {
			block, _ := pem.Decode(caKey)
			caKey = block.Bytes
			klog.Info("Succeed in loading CA key from local directory")
		}

		if ca != nil && caKey != nil {
			Config.Ca = ca
			Config.CaKey = caKey
		} else if !(ca == nil && caKey == nil) {
			klog.Exit("Both of ca and caKey should be specified!")
		}

		cert, err := os.ReadFile(hub.TLSCertFile)
		if err == nil {
			block, _ := pem.Decode(cert)
			cert = block.Bytes
			klog.Info("Succeed in loading certificate from local directory")
		}
		key, err := os.ReadFile(hub.TLSPrivateKeyFile)
		if err == nil {
			block, _ := pem.Decode(key)
			key = block.Bytes
			klog.Info("Succeed in loading private key from local directory")
		}

		if cert != nil && key != nil {
			Config.Cert = cert
			Config.Key = key
		} else if !(cert == nil && key == nil) {
			klog.Exit("Both of cert and key should be specified!")
		}
	})
}

// ObjectSyncController use beehive context message layer
type ObjectSyncController struct {
	CrdClient versioned.Interface

	// informer
	ClusterObjectSyncInformer syncinformer.ClusterObjectSyncInformer
	ObjectSyncInformer        syncinformer.ObjectSyncInformer

	// synced
	ClusterObjectSyncSynced cache.InformerSynced
	ObjectSyncSynced        cache.InformerSynced

	// lister
	ClusterObjectSyncLister synclister.ClusterObjectSyncLister
	ObjectSyncLister        synclister.ObjectSyncLister
}
