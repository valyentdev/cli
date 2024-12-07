package machines

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	machinesCmd := &cobra.Command{
		Use: "machines",
	}
	machinesCmd.AddCommand(newListMachinesCmd())
	machinesCmd.AddCommand(newLogsCmd())
	machinesCmd.AddCommand(newListMachineEventsCmd())
	machinesCmd.AddCommand(newStartMachineCmd())
	machinesCmd.AddCommand(newStopMachineCmd())

	return machinesCmd
}
