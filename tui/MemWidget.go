package tui

import (
	"time"

	"github.com/gizak/termui/v3/widgets"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/memory"
)

type MemWidget struct {
	*widgets.Plot
	updateInterval time.Duration
	memory         *memory.Memory
}

// type MemoryInfo struct {
// 	Total       uint64
// 	Used        uint64
// 	UsedPercent float64
// }

// func (self *MemWidget) renderMemInfo(line string, memoryInfo MemoryInfo) {
// 	self.Data[line] = append(self.Data[line], memoryInfo.UsedPercent)
// 	memoryTotalBytes, memoryTotalMagnitude := utils.ConvertBytes(memoryInfo.Total)
// 	memoryUsedBytes, memoryUsedMagnitude := utils.ConvertBytes(memoryInfo.Used)
// 	self.Labels[line] = fmt.Sprintf("%3.0f%% %5.1f%s/%.0f%s",
// 		memoryInfo.UsedPercent,
// 		memoryUsedBytes,
// 		memoryUsedMagnitude,
// 		memoryTotalBytes,
// 		memoryTotalMagnitude,
// 	)
// }

func (self *MemWidget) updateMainMemory() {
	// self.renderMemInfo("Main", MemoryInfo{
	// 	Total:       mainMemory.Total,
	// 	Used:        mainMemory.Used,
	// 	UsedPercent: mainMemory.UsedPercent,
	// })
	self.Data[0] = append([]float64{float64(len(self.memory.PhysicalMemory))}, self.Data[0]...)
}

func (self *MemWidget) updateVirtualMemory() {
	// self.renderMemInfo("Virtual", MemoryInfo{
	// 	Total:       mainMemory.Total,
	// 	Used:        mainMemory.Used,
	// 	UsedPercent: mainMemory.UsedPercent,
	// })

	self.Data[1] = append([]float64{float64(len(self.memory.VirtualMemory))}, self.Data[1]...)

}

func NewMemWidget(mem *memory.Memory) *MemWidget {
	self := &MemWidget{
		Plot:           widgets.NewPlot(),
		updateInterval: time.Second,
		memory:         mem,
	}
	self.Title = " Memory Usage "
	self.HorizontalScale = 7
	self.DrawDirection = widgets.DrawRight

	self.Data = make([][]float64, 2)
	self.Data[0] = []float64{0}
	self.Data[1] = []float64{0}

	self.updateMainMemory()
	self.updateVirtualMemory()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			self.updateMainMemory()
			self.updateVirtualMemory()
			self.Unlock()
		}
	}()

	return self
}
