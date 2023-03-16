package widget

import (
	"fmt"
	"strings"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mat"
	"github.com/mokiat/lacking/util/metrics"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const (
	regionHeight = 45
)

type RegionBlockData struct {
	Regions []metrics.RegionStat
}

var RegionBlock = co.Define(func(props co.Properties, scope co.Scope) co.Instance {
	var (
		data = co.GetData[RegionBlockData](props)
	)

	var maxDepth int
	for _, region := range data.Regions {
		if region.Depth > maxDepth {
			maxDepth = region.Depth
		}
	}

	essence := co.UseState(func() *regionBlockEssence {
		return &regionBlockEssence{
			selectedNodeID: metrics.NilParentID,
			placement:      make(map[int]regionPlacement),
			graph:          make(map[int]regionNode),
		}
	}).Get()
	essence.font = co.OpenFont(scope, "mat:///roboto-bold.ttf")
	essence.rebuildGraph(data.Regions)

	return co.New(mat.Element, func() {
		co.WithData(mat.ElementData{
			Essence: essence,
			IdealSize: opt.V(ui.Size{
				Width:  200,
				Height: maxDepth * regionHeight,
			}),
		})
		co.WithLayoutData(props.LayoutData())
	})
})

var _ ui.ElementMouseHandler = (*regionBlockEssence)(nil)
var _ ui.ElementRenderHandler = (*regionBlockEssence)(nil)

type regionBlockEssence struct {
	font     *ui.Font
	fontSize float32

	selectedNodeID int
	graph          map[int]regionNode
	placement      map[int]regionPlacement
}

func (b *regionBlockEssence) rebuildGraph(stats []metrics.RegionStat) {
	slices.SortFunc(stats, func(a, b metrics.RegionStat) bool {
		if a.Depth == b.Depth {
			return a.ID > b.ID // inversed on purpose
		}
		return a.Depth < b.Depth
	})

	rootDuration := time.Duration(0)
	for _, stat := range stats {
		if stat.ParentID == metrics.NilParentID {
			rootDuration += stat.Duration
		}
	}

	rootSamples := 0
	for _, stat := range stats {
		if stat.ParentID == metrics.NilParentID {
			if stat.Samples > rootSamples {
				rootSamples = stat.Samples
			}
		}
	}

	maps.Clear(b.graph)
	if rootSamples == 0 {
		return
	}
	b.graph[metrics.NilParentID] = regionNode{
		ParentID:      metrics.NilParentID,
		ID:            metrics.NilParentID,
		FirstChild:    metrics.NilParentID,
		Duration:      rootDuration / time.Duration(rootSamples),
		DurationRatio: 1.0,
	}
	for _, stat := range stats {
		parentNode := b.graph[stat.ParentID]
		node := regionNode{
			ParentID:      stat.ParentID,
			ID:            stat.ID,
			NextSibling:   parentNode.FirstChild,
			FirstChild:    metrics.NilParentID,
			Name:          stripRegionNamespace(stat.Name),
			Duration:      stat.Duration / time.Duration(rootSamples),
			DurationRatio: float32(stat.Duration.Seconds() / rootDuration.Seconds()),
		}
		parentNode.FirstChild = stat.ID
		b.graph[stat.ParentID] = parentNode
		b.graph[stat.ID] = node
	}
}

func (b *regionBlockEssence) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	if event.Type == ui.MouseEventTypeDown {
		if event.Button == ui.MouseButtonLeft {
			position := event.Position.Translate(element.ContentBounds().Position.Inverse())
			for id, placement := range b.placement {
				bounds := ui.NewBounds(
					int(placement.Left),
					int(placement.Top),
					int(placement.Width),
					int(placement.Height),
				)
				if bounds.Contains(position) {
					b.selectedNodeID = id
					return true
				}
			}
		}
		if event.Button == ui.MouseButtonRight {
			node := b.graph[b.selectedNodeID]
			b.selectedNodeID = node.ParentID
			return true
		}
	}
	return false
}

func (b *regionBlockEssence) OnRender(element *ui.Element, canvas *ui.Canvas) {
	node := b.graph[b.selectedNodeID]
	maps.Clear(b.placement)
	b.renderRegionNode(element, canvas, node)
}

func (b *regionBlockEssence) renderRegionNode(element *ui.Element, canvas *ui.Canvas, node regionNode) {
	if node.Duration == 0 {
		return
	}

	parentDuration := b.graph[node.ParentID].Duration
	parentPlacement, ok := b.placement[node.ParentID]
	if !ok {
		contentArea := element.ContentBounds()
		parentPlacement = regionPlacement{
			Top:      0.0,
			Left:     0.0,
			Height:   0.0,
			Width:    float32(contentArea.Width),
			FreeLeft: 0.0,
		}
		parentDuration = node.Duration
	}

	if node.ID != metrics.NilParentID {
		regionPosition := sprec.Vec2{
			X: parentPlacement.FreeLeft,
			Y: parentPlacement.Top + parentPlacement.Height,
		}
		regionSize := sprec.Vec2{
			X: parentPlacement.Width * float32(node.Duration.Seconds()/parentDuration.Seconds()),
			Y: regionHeight,
		}
		parentPlacement.FreeLeft += regionSize.X
		b.placement[node.ParentID] = parentPlacement

		b.placement[node.ID] = regionPlacement{
			Top:      regionPosition.Y,
			Left:     regionPosition.X,
			Height:   regionSize.Y,
			Width:    regionSize.X,
			FreeLeft: regionPosition.X,
		}

		canvas.SetClipRect(
			regionPosition.X,
			regionPosition.X+regionSize.X,
			regionPosition.Y,
			regionPosition.Y+regionSize.Y,
		)
		canvas.Reset()
		canvas.SetStrokeColor(ui.Gray())
		canvas.SetStrokeSizeSeparate(1.0, 0.0)
		canvas.Rectangle(regionPosition, regionSize)
		canvas.Fill(ui.Fill{
			Color: ui.ColorWithAlpha(ui.Navy(), 196),
		})
		canvas.Stroke()

		name := node.Name
		duration := node.Duration.Truncate(10 * time.Microsecond)
		percentage := node.DurationRatio * 100
		text := fmt.Sprintf("%s\n%s | %.2f%%", name, duration, percentage)
		fontSize := float32(20)
		textSize := b.font.TextSize(text, fontSize)
		textPosition := sprec.Vec2{
			X: regionPosition.X + (regionSize.X-textSize.X)/2.0,
			Y: regionPosition.Y + (regionSize.Y-textSize.Y)/2.0,
		}
		canvas.FillText(text, textPosition, ui.Typography{
			Font:  b.font,
			Size:  fontSize,
			Color: ui.White(),
		})
	}

	childID := node.FirstChild
	for childID != metrics.NilParentID {
		childNode := b.graph[childID]
		b.renderRegionNode(element, canvas, childNode)
		childID = childNode.NextSibling
	}
}

func stripRegionNamespace(name string) string {
	if index := strings.IndexRune(name, ':'); index >= 0 {
		return name[index+1:]
	}
	return name
}

type regionNode struct {
	ParentID      int
	ID            int
	NextSibling   int
	FirstChild    int
	Name          string
	Duration      time.Duration
	DurationRatio float32
}

type regionPlacement struct {
	Top      float32
	Left     float32
	Height   float32
	Width    float32
	FreeLeft float32
}
