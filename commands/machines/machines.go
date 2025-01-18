package machines

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	machinesCmd := &cobra.Command{
		Use:     "machines",
		Short:   "m",
		Aliases: []string{"machine"},
	}
	machinesCmd.AddCommand(newListMachinesCmd())
	machinesCmd.AddCommand(newLogsCmd())
	machinesCmd.AddCommand(newListMachineEventsCmd())
	machinesCmd.AddCommand(newStartMachineCmd())
	machinesCmd.AddCommand(newStopMachineCmd())
	machinesCmd.AddCommand(newDeleteMachineCmd())
	machinesCmd.AddCommand(newCreateMachineCmd())
	machinesCmd.AddCommand(newExecMachineCmd())

	return machinesCmd
}
