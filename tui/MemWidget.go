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

func (m *MemWidget) updateMainMemory() {
	// m.renderMemInfo("Main", MemoryInfo{
	// 	Total:       mainMemory.Total,
	// 	Used:        mainMemory.Used,
	// 	UsedPercent: mainMemory.UsedPercent,
	// })
	m.Data[0] = append([]float64{float64(len(m.memory.PhysicalMemory))}, m.Data[0]...)
}

func (m *MemWidget) updateVirtualMemory() {
	// m.renderMemInfo("Virtual", MemoryInfo{
	// 	Total:       mainMemory.Total,
	// 	Used:        mainMemory.Used,
	// 	UsedPercent: mainMemory.UsedPercent,
	// })

	m.Data[1] = append([]float64{float64(len(m.memory.VirtualMemory))}, m.Data[1]...)

}

func NewMemWidget(mem *memory.Memory) *MemWidget {
	m := &MemWidget{
		Plot:           widgets.NewPlot(),
		updateInterval: time.Second,
		memory:         mem,
	}
	m.Title = " Memory Usage "
	m.HorizontalScale = 7
	m.DrawDirection = widgets.DrawRight

	m.Data = make([][]float64, 2)
	m.Data[0] = []float64{0}
	m.Data[1] = []float64{0}

	m.updateMainMemory()
	m.updateVirtualMemory()

	go func() {
		for range time.NewTicker(m.updateInterval).C {
			m.Lock()
			m.updateMainMemory()
			m.updateVirtualMemory()
			m.Unlock()
		}
	}()

	return m
}
