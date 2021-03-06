package commands_test

import (
	"errors"

	"github.com/cloudfoundry/bosh-bootloader/commands"
	"github.com/cloudfoundry/bosh-bootloader/fakes"
	"github.com/cloudfoundry/bosh-bootloader/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("create-lbs", func() {
	var (
		command        commands.CreateLBs
		awsCreateLBs   *fakes.AWSCreateLBs
		gcpCreateLBs   *fakes.GCPCreateLBs
		stateValidator *fakes.StateValidator
	)

	BeforeEach(func() {
		awsCreateLBs = &fakes.AWSCreateLBs{}
		gcpCreateLBs = &fakes.GCPCreateLBs{}
		stateValidator = &fakes.StateValidator{}

		command = commands.NewCreateLBs(awsCreateLBs, gcpCreateLBs, stateValidator)
	})

	Describe("Execute", func() {
		It("creates a GCP lb type if the iaas if GCP", func() {
			err := command.Execute([]string{
				"--type", "concourse",
				"--skip-if-exists",
			}, storage.State{
				IAAS: "gcp",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(gcpCreateLBs.ExecuteCall.Receives.Config).Should(Equal(commands.GCPCreateLBsConfig{
				LBType:       "concourse",
				SkipIfExists: true,
			}))
		})

		It("creates a GCP cf lb type is the iaas if GCP and type is cf", func() {
			err := command.Execute([]string{
				"--type", "cf",
				"--cert", "my-cert",
				"--key", "my-key",
				"--skip-if-exists",
			}, storage.State{
				IAAS: "gcp",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(gcpCreateLBs.ExecuteCall.Receives.Config).Should(Equal(commands.GCPCreateLBsConfig{
				LBType:       "cf",
				CertPath:     "my-cert",
				KeyPath:      "my-key",
				SkipIfExists: true,
			}))
		})

		It("creates an AWS lb type if the iaas is AWS", func() {
			err := command.Execute([]string{
				"--type", "concourse",
				"--cert", "my-cert",
				"--key", "my-key",
				"--chain", "my-chain",
				"--skip-if-exists", "true",
			}, storage.State{
				IAAS: "aws",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(awsCreateLBs.ExecuteCall.Receives.Config).Should(Equal(commands.AWSCreateLBsConfig{
				LBType:       "concourse",
				CertPath:     "my-cert",
				KeyPath:      "my-key",
				ChainPath:    "my-chain",
				SkipIfExists: true,
			}))
		})

		Context("failure cases", func() {
			It("returns an error when state validator fails", func() {
				stateValidator.ValidateCall.Returns.Error = errors.New("state validator failed")
				err := command.Execute([]string{}, storage.State{})

				Expect(stateValidator.ValidateCall.CallCount).To(Equal(1))
				Expect(err).To(MatchError("state validator failed"))
			})

			It("returns an error when an invalid command line flag is supplied", func() {
				err := command.Execute([]string{"--invalid-flag"}, storage.State{})
				Expect(err).To(MatchError("flag provided but not defined: -invalid-flag"))
			})

			It("returns an error when the AWSCreateLBs fails", func() {
				awsCreateLBs.ExecuteCall.Returns.Error = errors.New("something bad happened")

				err := command.Execute([]string{"some-aws-args"}, storage.State{
					IAAS: "aws",
				})
				Expect(err).To(MatchError("something bad happened"))
			})

			It("returns an error when the GCPCreateLBs fails", func() {
				gcpCreateLBs.ExecuteCall.Returns.Error = errors.New("something bad happened")

				err := command.Execute([]string{"some-gcp-args"}, storage.State{
					IAAS: "gcp",
				})
				Expect(err).To(MatchError("something bad happened"))
			})
		})
	})
})
