package widget

import (
	"cmp"
	"fmt"
	"strings"
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
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

var RegionBlock = co.Define(&regionBlockComponent{})

type regionBlockComponent struct {
	co.BaseComponent

	font *ui.Font

	maxDepth       int
	selectedNodeID int
	graph          map[int]regionNode
	placement      map[int]regionPlacement
}

func (c *regionBlockComponent) OnCreate() {
	c.font = co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf")
	c.selectedNodeID = metrics.NilParentID
	c.placement = make(map[int]regionPlacement)
	c.graph = make(map[int]regionNode)
}

func (c *regionBlockComponent) OnUpsert() {
	data := co.GetData[RegionBlockData](c.Properties())

	c.maxDepth = 0
	for _, region := range data.Regions {
		if region.Depth > c.maxDepth {
			c.maxDepth = region.Depth
		}
	}
	c.rebuildGraph(data.Regions)
}

func (c *regionBlockComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence: c,
			IdealSize: opt.V(ui.Size{
				Width:  200,
				Height: c.maxDepth * regionHeight,
			}),
		})
	})
}

func (c *regionBlockComponent) rebuildGraph(stats []metrics.RegionStat) {
	slices.SortFunc(stats, func(a, b metrics.RegionStat) int {
		if a.Depth == b.Depth {
			return cmp.Compare(a.ID, b.ID)
		}
		return cmp.Compare(a.Depth, b.Depth)
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

	maps.Clear(c.graph)
	if rootSamples == 0 {
		return
	}
	c.graph[metrics.NilParentID] = regionNode{
		ParentID:      metrics.NilParentID,
		ID:            metrics.NilParentID,
		FirstChild:    metrics.NilParentID,
		Duration:      rootDuration / time.Duration(rootSamples),
		DurationRatio: 1.0,
	}
	for _, stat := range stats {
		parentNode := c.graph[stat.ParentID]
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
		c.graph[stat.ParentID] = parentNode
		c.graph[stat.ID] = node
	}
}

func (c *regionBlockComponent) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	if event.Action == ui.MouseActionDown {
		if event.Button == ui.MouseButtonLeft {
			position := event.Position().Translate(element.ContentBounds().Position.Inverse())
			for id, placement := range c.placement {
				bounds := ui.NewBounds(
					int(placement.Left),
					int(placement.Top),
					int(placement.Width),
					int(placement.Height),
				)
				if bounds.Contains(position) {
					c.selectedNodeID = id
					return true
				}
			}
		}
		if event.Button == ui.MouseButtonRight {
			node := c.graph[c.selectedNodeID]
			c.selectedNodeID = node.ParentID
			return true
		}
	}
	return false
}

func (c *regionBlockComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	node := c.graph[c.selectedNodeID]
	maps.Clear(c.placement)
	c.renderRegionNode(element, canvas, node)
}

func (c *regionBlockComponent) renderRegionNode(element *ui.Element, canvas *ui.Canvas, node regionNode) {
	if node.Duration == 0 {
		return
	}

	parentDuration := c.graph[node.ParentID].Duration
	parentPlacement, ok := c.placement[node.ParentID]
	if !ok {
		drawBounds := canvas.DrawBounds(element, false)
		parentPlacement = regionPlacement{
			Top:      0.0,
			Left:     0.0,
			Height:   0.0,
			Width:    drawBounds.Width(),
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
		c.placement[node.ParentID] = parentPlacement

		c.placement[node.ID] = regionPlacement{
			Top:      regionPosition.Y,
			Left:     regionPosition.X,
			Height:   regionSize.Y,
			Width:    regionSize.X,
			FreeLeft: regionPosition.X,
		}

		canvas.ClipRect(
			regionPosition,
			regionSize,
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
		textSize := c.font.TextSize(text, fontSize)
		textPosition := sprec.Vec2{
			X: regionPosition.X + (regionSize.X-textSize.X)/2.0,
			Y: regionPosition.Y + (regionSize.Y-textSize.Y)/2.0,
		}
		canvas.FillTextLine([]rune(text), textPosition, ui.Typography{
			Font:  c.font,
			Size:  fontSize,
			Color: ui.White(),
		})
	}

	childID := node.FirstChild
	for childID != metrics.NilParentID {
		childNode := c.graph[childID]
		c.renderRegionNode(element, canvas, childNode)
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
