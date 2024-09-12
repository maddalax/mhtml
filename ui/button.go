package ui

import (
	"mhtml/h"
)

type ButtonProps struct {
	Id       string
	Text     string
	Target   string
	Type     string
	Trigger  string
	Get      string
	Class    string
	Children []h.Renderable
}

func PrimaryButton(props ButtonProps) h.Renderable {
	props.Class = h.MergeClasses(props.Class, "border-blue-700 bg-blue-700 text-white")
	return Button(props)
}

func SecondaryButton(props ButtonProps) h.Renderable {
	props.Class = h.MergeClasses(props.Class, "border-gray-700 bg-gray-700 text-white")
	return Button(props)
}

func Button(props ButtonProps) h.Renderable {

	text := h.Text(props.Text)

	button := h.Button(
		h.If(props.Id != "", h.Id(props.Id)),
		h.If(props.Children != nil, h.Children(props.Children...)),
		h.If(props.Trigger != "", h.Trigger(props.Trigger)),
		h.Class("flex gap-1 items-center border p-4 rounded cursor-hover", props.Class),
		h.If(props.Get != "", h.Get(props.Get)),
		h.If(props.Target != "", h.Target(props.Target)),
		//h.Attribute("hx-indicator", "#spinner"),
		h.IfElse(props.Type != "", h.Type(props.Type), h.Type("button")),
		h.BeforeRequestSetText("Loading..."),
		h.AfterRequestSetText(props.Text),
		text,
	)

	return button
}
