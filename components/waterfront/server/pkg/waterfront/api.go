/*
Copyright 2018 Pax Automa Systems, Inc.

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

package waterfront

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kube_v1 "k8s.io/client-go/pkg/api/v1"

	teamster_proto "github.com/paxautoma/operos/components/teamster/pkg/teamster"
)

type WaterfrontAPI struct {
	teamsterClient teamster_proto.TeamsterClient
	kubeClient     *kubernetes.Clientset
}

func NewWaterfrontAPI(teamsterClient teamster_proto.TeamsterClient, kubeClient *kubernetes.Clientset) *WaterfrontAPI {
	return &WaterfrontAPI{
		teamsterClient: teamsterClient,
		kubeClient:     kubeClient,
	}
}

func (w *WaterfrontAPI) ListNodes(ctx context.Context, empty *Empty) (*ListNodesResponse, error) {
	res, err := w.teamsterClient.ListNodes(ctx, &teamster_proto.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "error accessing teamster")
	}

	kubeNodeList, err := w.kubeClient.Nodes().List(meta_v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed fetching node list from kube")
	}

	kubeNodeMap := make(map[string]*kube_v1.Node)
	for _, kubeNode := range kubeNodeList.Items {
		kubeNode := kubeNode
		kubeNodeMap[kubeNode.GetName()] = &kubeNode
	}

	nodes := make([]*Node, len(res.Nodes))
	teamsterNodeIds := make(map[string]bool)
	for idx, node := range res.Nodes {
		kubeNode, _ := kubeNodeMap[node.Uuid]
		nodes[idx] = nodeFromKube(node.Uuid, kubeNode)
		teamsterNodeIds[node.Uuid] = true
	}

	for _, kubeNode := range kubeNodeList.Items {
		nodeID := kubeNode.GetName()
		if _, inTeamster := teamsterNodeIds[nodeID]; !inTeamster {
			nodes = append(nodes, nodeFromKube(nodeID, &kubeNode))
		}
	}

	return &ListNodesResponse{Nodes: nodes}, nil
}

func nodeFromKube(uuid string, kubeNode *kube_v1.Node) *Node {
	node := &Node{
		Id:     uuid,
		Status: NodeStatus_NOT_READY,
	}

	if kubeNode != nil {
		node.PodCidr = kubeNode.Spec.PodCIDR
		node.Ip = kubeNode.Status.Addresses[0].Address
		node.Status = nodeReady(kubeNode)
	}

	return node
}

func nodeReady(kubeNode *kube_v1.Node) NodeStatus {
	for _, condition := range kubeNode.Status.Conditions {
		if condition.Type == kube_v1.NodeReady {
			if condition.Status == "True" {
				return NodeStatus_READY
			}
			return NodeStatus_NOT_READY
		}
	}
	return NodeStatus_NOT_READY
}

func (w *WaterfrontAPI) GetNode(ctx context.Context, req *GetNodeRequest) (*GetNodeResponse, error) {
	res, err := w.teamsterClient.GetNodeHardware(ctx, &teamster_proto.GetNodeHardwareRequest{Uuid: req.Id})
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return nil, errors.Wrap(err, "error accessing teamster")
		}
	}

	kubeNode, err := w.kubeClient.Nodes().Get(req.Id, meta_v1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error fetching node info from kube")
	}

	node := nodeFromKube(req.Id, kubeNode)
	if res != nil {
		node.HardwareInfo = res.HardwareInfo
	}

	return &GetNodeResponse{Node: node}, nil
}

func (w *WaterfrontAPI) readSettingsFile() (map[string]string, error) {
	fp, err := os.Open("/etc/paxautoma/settings")
	if err != nil {
		return nil, errors.Wrap(err, "could not open settings file")
	}

	result := make(map[string]string)

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, errors.Wrap(err, "failed to read settings file")
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			return nil, errors.Errorf("could not parse this line in settings file: '%s'", line)
		}

		result[parts[0]] = strings.Trim(parts[1], "\"")
	}

	return result, nil
}

func (w *WaterfrontAPI) GetClusterInfo(ctx context.Context, req *Empty) (*GetClusterInfoResponse, error) {
	res, err := w.teamsterClient.GetCACertExpiry(ctx, &teamster_proto.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "error accessing teamster")
	}

	settings, err := w.readSettingsFile()
	if err != nil {
		return nil, err
	}

	return &GetClusterInfoResponse{
		LicenseExpiry: res.ExpiryUnix,
		Settings:      settings,
	}, nil
}

func (w *WaterfrontAPI) SetRootPassword(ctx context.Context, req *SetRootPasswordRequest) (*Empty, error) {
	_, err := w.teamsterClient.SetRootPassword(ctx, &teamster_proto.SetRootPasswordRequest{req.Password})
	if err != nil {
		return nil, errors.Wrap(err, "error accessing teamster")
	}

	return &Empty{}, nil
}
