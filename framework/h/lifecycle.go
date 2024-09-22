package h

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/hx"
	"github.com/maddalax/htmgo/framework/internal/util"
)

type LifeCycle struct {
	handlers map[hx.Event][]Command
}

func NewLifeCycle() *LifeCycle {
	return &LifeCycle{
		handlers: make(map[hx.Event][]Command),
	}
}

func validateCommands(cmds []Command) {
	for _, cmd := range cmds {
		switch t := cmd.(type) {
		case SimpleJsCommand:
			break
		case ComplexJsCommand:
			break
		case *AttributeMap:
			break
		case *Element:
			panic(fmt.Sprintf("element is not allowed in lifecycle events. Got: %v", t))
		default:
			panic(fmt.Sprintf("type is not allowed in lifecycle events. Got: %v", t))

		}
	}
}

func (l *LifeCycle) OnEvent(event hx.Event, cmd ...Command) *LifeCycle {
	validateCommands(cmd)

	if l.handlers[event] == nil {
		l.handlers[event] = []Command{}
	}

	l.handlers[event] = append(l.handlers[event], cmd...)
	return l
}

func (l *LifeCycle) BeforeRequest(cmd ...Command) *LifeCycle {
	l.OnEvent(hx.BeforeRequestEvent, cmd...)
	return l
}

func OnLoad(cmd ...Command) *LifeCycle {
	return NewLifeCycle().OnEvent(hx.LoadEvent, cmd...)
}

func OnAfterSwap(cmd ...Command) *LifeCycle {
	return NewLifeCycle().OnEvent(hx.AfterSwapEvent, cmd...)
}

func OnTrigger(trigger string, cmd ...Command) *LifeCycle {
	return NewLifeCycle().OnEvent(hx.NewStringTrigger(trigger).ToString(), cmd...)
}

func OnClick(cmd ...Command) *LifeCycle {
	return NewLifeCycle().OnEvent(hx.ClickEvent, cmd...)
}

func OnEvent(event hx.Event, cmd ...Command) *LifeCycle {
	return NewLifeCycle().OnEvent(event, cmd...)
}

func BeforeRequest(cmd ...Command) *LifeCycle {
	return NewLifeCycle().BeforeRequest(cmd...)
}

func AfterRequest(cmd ...Command) *LifeCycle {
	return NewLifeCycle().AfterRequest(cmd...)
}

func OnMutationError(cmd ...Command) *LifeCycle {
	return NewLifeCycle().OnMutationError(cmd...)
}

func (l *LifeCycle) AfterRequest(cmd ...Command) *LifeCycle {
	l.OnEvent(hx.AfterRequestEvent, cmd...)
	return l
}

func (l *LifeCycle) OnMutationError(cmd ...Command) *LifeCycle {
	l.OnEvent(hx.OnMutationErrorEvent, cmd...)
	return l
}

type Command = Ren

type SimpleJsCommand struct {
	Command string
}

type ComplexJsCommand struct {
	Command      string
	TempFuncName string
}

func SetText(text string) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("this.innerText = '%s'", text)}
}

func Increment(amount int) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("this.innerText = parseInt(this.innerText) + %d", amount)}
}

func SetInnerHtml(r Ren) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("this.innerHTML = `%s`", Render(r))}
}

func SetOuterHtml(r Ren) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("this.outerHTML = `%s`", Render(r))}
}

func AddAttribute(name, value string) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("this.setAttribute('%s', '%s')", name, value)}
}

func SetDisabled(disabled bool) SimpleJsCommand {
	if disabled {
		return AddAttribute("disabled", "true")
	} else {
		return RemoveAttribute("disabled")
	}
}

func RemoveAttribute(name string) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("this.removeAttribute('%s')", name)}
}

func AddClass(class string) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("this.classList.add('%s')", class)}
}

func RemoveClass(class string) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("this.classList.remove('%s')", class)}
}

func ToggleClass(class string) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("this.classList.toggle('%s')", class)}
}

func ToggleClassOnElement(selector, class string) ComplexJsCommand {
	// language=JavaScript
	return EvalJs(fmt.Sprintf(`
		var el = document.querySelector('%s');
		if(el) { el.classList.toggle('%s'); }`,
	))
}

func Alert(text string) SimpleJsCommand {
	// language=JavaScript
	return SimpleJsCommand{Command: fmt.Sprintf("alert('%s')", text)}
}

func EvalJs(js string) ComplexJsCommand {
	name := fmt.Sprintf("__eval_%s", util.RandSeq(6))
	return ComplexJsCommand{Command: js, TempFuncName: name}
}

func InjectScript(src string) ComplexJsCommand {
	// language=JavaScript
	return ComplexJsCommand{Command: fmt.Sprintf(`
		var script = document.createElement('script');
		script.src = '%s';
        src.async = true;
		document.head.appendChild(script);
	`, src)}
}

func InjectScriptIfNotExist(src string) ComplexJsCommand {
	// language=JavaScript
	return EvalJs(fmt.Sprintf(`
		if(!document.querySelector('script[src="%s"]')) {
			var script = document.createElement('script');
			script.src = '%s';
			script.async = true;
			document.head.appendChild(script);
		}
	`, src, src))
}
