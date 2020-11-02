package main

import (
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/core"
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/network"
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/servicebus"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create an Azure Resource Group
		tags := pulumi.StringMap{
			"app":     pulumi.String("azuretest"),
			"version": pulumi.String("1.0"),
		}

		resourceGroup, err := core.NewResourceGroup(ctx, "amandasouza-rg", &core.ResourceGroupArgs{
			Location: pulumi.String("WestUS"),
			Tags:     tags,
		})
		if err != nil {
			return err
		}

		// Create a Virtual Network.
		vnetArgs := network.VirtualNetworkArgs{
			ResourceGroupName: resourceGroup.Name,
			Location:          resourceGroup.Location,
			AddressSpaces: pulumi.StringArray{
				pulumi.String("10.2.0.0/16"),
			},
		}
		vnet, err := network.NewVirtualNetwork(ctx, "azuretest-vnet", &vnetArgs)
		if err != nil {
			return err
		}

		// Create a subnet.
		subnetArgs := network.SubnetArgs{
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: vnet.Name,
			AddressPrefixes: pulumi.StringArray{
				pulumi.String("10.2.1.0/24"),
			},
		}
		_, err = network.NewSubnet(ctx, "azuretest-subnet", &subnetArgs)
		if err != nil {
			return err
		}

		//createServiceBus
		namespace, err := servicebus.NewNamespace(ctx, "azuretest-pipeline", &servicebus.NamespaceArgs{
			Location:          resourceGroup.Location,
			ResourceGroupName: resourceGroup.Name,
			Sku:               pulumi.String("Standard"),
			Tags:              tags,
		})
		if err != nil {
			return err
		}

		_, err = servicebus.NewNamespaceAuthorizationRule(ctx, "azuretestAR", &servicebus.NamespaceAuthorizationRuleArgs{
			NamespaceName:     namespace.Name,
			ResourceGroupName: resourceGroup.Name,
			Listen:            pulumi.Bool(true),
			Send:              pulumi.Bool(true),
			Manage:            pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// _, err = servicebus.NewNamespaceNetworkRuleSet(ctx, "azuretestNR", &servicebus.NamespaceNetworkRuleSetArgs{
		// 	NamespaceName:     namespace.Name,
		// 	ResourceGroupName: resourceGroup.Name,
		// 	DefaultAction:     pulumi.String("Deny"),
		// 	NetworkRules: servicebus.NamespaceNetworkRuleSetNetworkRuleArray{
		// 		&servicebus.NamespaceNetworkRuleSetNetworkRuleArgs{
		// 			SubnetId:                         subnet.ID(),
		// 			IgnoreMissingVnetServiceEndpoint: pulumi.Bool(false),
		// 		},
		// 	},
		// 	IpRules: pulumi.StringArray{
		// 		pulumi.String("1.1.1.1"),
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }

		//createqueue
		_, err = servicebus.NewQueue(ctx, "azuretest-queue", &servicebus.QueueArgs{
			ResourceGroupName:  resourceGroup.Name,
			NamespaceName:      namespace.Name,
			EnablePartitioning: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		return nil

	})
}
