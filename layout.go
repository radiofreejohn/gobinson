package main

type LayoutNode struct {
	style_node StyledNode
	dimensions Dimensions
	children   []LayoutNode
}

type Dimensions struct {
	// position of the content area relative to the document origin
	x, y float32
	// content area size
	width, height float32
	// surrounding edges
	padding, border, margin EdgeSizes
}

type EdgeSizes struct {
	left, right, top, bottom float32
}

func layout(node StyledNode, containing_block Dimensions) LayoutNode {
	layout_node := LayoutNode{
		style_node: node,
		children:   make([]LayoutNode, 0)}

	calculate_width(&layout_node, containing_block) // may need to use pointer

	calculate_height(&layout_node, containing_block) // ^ that

	return layout_node
}

// this is just an idea
type BoolMatch struct {
	i, j, k bool
}

func calculate_width(node *LayoutNode, containing_block Dimensions) {
	style := node.style_node

	auto := KeywordValue("auto")

	width, ok := style.value("width") // style.value("width").unwrap_or(auto.clone());
	if !ok {
		width = auto
	}

	zero := LengthValue{0.0, Px}

	// rust passes address of zero here, but it's not modified, so I am
	// not passing address because stuff
	margin_left := style.lookup("margin-left", "margin", zero)
	margin_right := style.lookup("margin-right", "margin", zero)

	border_left := style.lookup("border-left", "border", zero)
	border_right := style.lookup("border-right", "border", zero)

	padding_left := style.lookup("padding-left", "padding", zero)
	padding_right := style.lookup("padding-right", "padding", zero)

	total := 0.0
	for _, val := range []Value{margin_left, margin_right, border_left, border_right, padding_left, padding_right} {
		total += val.to_px()
	}

	// If width is not auto and the total is wider than the container, treat auto margins as 0.
	if width != auto && total > containing_block.width {
		if m, ok := margin_left.(KeywordValue); ok && m == KeywordValue("auto") {
			margin_left = LengthValue{0.0, Px}
		}
		if m, ok = margin_right.(KeywordValue); ok && m == KeywordValue("auto") {
			margin_right = LengthValue{0.0, Px}
		}
	}

	underflow := containing_block.width - total
	match := BoolMatch{width == auto, margin_left == auto, margin_right == auto}
	switch {
	// If the values are overconstrained, calculate margin_right.
	case match == BoolMatch{false, false, false}:
		margin_right = LengthValue{margin_right.to_px() + underflow, Px}
	// If exactly one value is auto, its used value follows from the equality.
	case match == BoolMatch{false, false, true}:
		margin_right = LengthValue{underflow, Px}
	case match == BoolMatch{false, true, false}:
		margin_left = LengthValue{underflow, Px}
	// If margin-left and margin-right are both auto, their used values are equal.
	case match == BoolMatch{false, true, true}:
		margin_left = LengthValue{underflow / 2.0, Px}
		margin_right = LengthValue{underflow / 2.0, Px}
	// If width is set to auto, any other auto values become 0.
	case match.i == true:
		if margin_left == auto {
			margin_left = LengthValue{0.0, Px}
		}
		if margin_right == auto {
			margin_right = LengthValue{0.0, Px}
		}
		width = LengthValue{underflow, Px}
	}

	d := &node.dimensions
	d.width = width.to_px()

	d.padding.left = padding_left.to_px()
	d.padding.right = padding_right.to_px()

	d.border.left = border_left.to_px()
	d.border.right = border_right.to_px()

	d.margin.left = margin_left.to_px()
	d.margin.right = margin_right.to_px()

	d.x = containing_block.x + d.margin.left + d.border.left + d.padding.left
}
