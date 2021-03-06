package bosh_test

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/cloudfoundry/bosh-bootloader/bosh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/gomegamatchers"
)

var _ = Describe("CloudConfigGenerator", func() {
	Describe("Generate", func() {
		var (
			cloudConfigGenerator bosh.CloudConfigGenerator
		)

		BeforeEach(func() {
			cloudConfigGenerator = bosh.NewCloudConfigGenerator()
		})

		It("returns a generated cloud config that matches our example fixture", func() {
			cloudConfig, err := cloudConfigGenerator.Generate(bosh.CloudConfigInput{
				AZs: []string{"us-east-1a", "us-east-1b", "us-east-1c"},
				Subnets: []bosh.SubnetInput{
					{
						AZ:     "us-east-1a",
						Subnet: "some-subnet-1",
						CIDR:   "10.0.16.0/20",
						SecurityGroups: []string{
							"some-security-group-1",
						},
					},
					{
						AZ:     "us-east-1b",
						Subnet: "some-subnet-2",
						CIDR:   "10.0.32.0/20",
						SecurityGroups: []string{
							"some-security-group-2",
						},
					},
					{
						AZ:     "us-east-1c",
						Subnet: "some-subnet-3",
						CIDR:   "10.0.48.0/20",
						SecurityGroups: []string{
							"some-security-group-3",
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			buf, err := ioutil.ReadFile("fixtures/cloud_config_without_load_balancers.yml")
			Expect(err).NotTo(HaveOccurred())

			output, err := yaml.Marshal(cloudConfig)
			Expect(err).NotTo(HaveOccurred())

			Expect(output).To(MatchYAML(string(buf)))
		})

		Context("vm extensions", func() {
			It("generates a cloud config with load balancer vm extensions", func() {
				cloudConfig, err := cloudConfigGenerator.Generate(bosh.CloudConfigInput{
					LBs: []bosh.LoadBalancerExtension{
						{
							Name:    "first-lb",
							ELBName: "some-lb-1",
						},
						{
							Name:    "second-lb",
							ELBName: "some-lb-2",
							SecurityGroups: []string{
								"some-security-group",
								"some-other-security-group",
							},
						},
					},
					AZs: []string{"us-east-1a", "us-east-1b", "us-east-1c"},
					Subnets: []bosh.SubnetInput{
						{
							AZ:     "us-east-1a",
							Subnet: "some-subnet-1",
							CIDR:   "10.0.16.0/20",
							SecurityGroups: []string{
								"some-security-group-1",
							},
						},
						{
							AZ:     "us-east-1b",
							Subnet: "some-subnet-2",
							CIDR:   "10.0.32.0/20",
							SecurityGroups: []string{
								"some-security-group-2",
							},
						},
						{
							AZ:     "us-east-1c",
							Subnet: "some-subnet-3",
							CIDR:   "10.0.48.0/20",
							SecurityGroups: []string{
								"some-security-group-3",
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				buf, err := ioutil.ReadFile("fixtures/cloud_config_with_load_balancers.yml")
				Expect(err).NotTo(HaveOccurred())

				output, err := yaml.Marshal(cloudConfig)
				Expect(err).NotTo(HaveOccurred())

				Expect(output).To(MatchYAML(string(buf)))
			})
		})

		Context("failure cases", func() {
			It("returns an error when it fails to generate networks for manifest", func() {
				_, err := cloudConfigGenerator.Generate(bosh.CloudConfigInput{
					AZs: []string{"us-east-1a"},
					Subnets: []bosh.SubnetInput{
						{
							AZ:     "us-east-1a",
							Subnet: "some-subnet-1",
							CIDR:   "some-bad-cidr-block",
							SecurityGroups: []string{
								"some-security-group-1",
							},
						},
					},
				})

				Expect(err).To(MatchError(ContainSubstring("cannot parse CIDR block")))
			})
		})
	})
})
