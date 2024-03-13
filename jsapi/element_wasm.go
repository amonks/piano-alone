package jsapi

import "syscall/js"

type Element js.Value


func (el Element) GetAccessKey() string {
	return js.Value(el).Get("accessKey").String()
}

func (el Element) SetAccessKey(val string) {
	js.Value(el).Set("accessKey", val)
}


func (el Element) GetAccessKeyLabel() string {
	return js.Value(el).Get("accessKeyLabel").String()
}

func (el Element) SetAccessKeyLabel(val string) {
	js.Value(el).Set("accessKeyLabel", val)
}


func (el Element) GetAutocapitalize() string {
	return js.Value(el).Get("autocapitalize").String()
}

func (el Element) SetAutocapitalize(val string) {
	js.Value(el).Set("autocapitalize", val)
}


func (el Element) GetDir() string {
	return js.Value(el).Get("dir").String()
}

func (el Element) SetDir(val string) {
	js.Value(el).Set("dir", val)
}


func (el Element) GetDraggable() bool {
	return js.Value(el).Get("draggable").Bool()
}

func (el Element) SetDraggable(val bool) {
	js.Value(el).Set("draggable", val)
}


func (el Element) GetHidden() bool {
	return js.Value(el).Get("hidden").Bool()
}

func (el Element) SetHidden(val bool) {
	js.Value(el).Set("hidden", val)
}


func (el Element) GetInert() bool {
	return js.Value(el).Get("inert").Bool()
}

func (el Element) SetInert(val bool) {
	js.Value(el).Set("inert", val)
}


func (el Element) GetInnerText() string {
	return js.Value(el).Get("innerText").String()
}

func (el Element) SetInnerText(val string) {
	js.Value(el).Set("innerText", val)
}


func (el Element) GetLang() string {
	return js.Value(el).Get("lang").String()
}

func (el Element) SetLang(val string) {
	js.Value(el).Set("lang", val)
}


func (el Element) GetOffsetHeight() float64 {
	return js.Value(el).Get("offsetHeight").Float()
}

func (el Element) SetOffsetHeight(val float64) {
	js.Value(el).Set("offsetHeight", val)
}


func (el Element) GetOffsetLeft() float64 {
	return js.Value(el).Get("offsetLeft").Float()
}

func (el Element) SetOffsetLeft(val float64) {
	js.Value(el).Set("offsetLeft", val)
}


func (el Element) GetOffsetParent() js.Value {
	return js.Value(el).Get("offsetParent")
}

func (el Element) SetOffsetParent(val any) {
	js.Value(el).Set("offsetParent", val)
}


func (el Element) GetOffsetTop() float64 {
	return js.Value(el).Get("offsetTop").Float()
}

func (el Element) SetOffsetTop(val float64) {
	js.Value(el).Set("offsetTop", val)
}


func (el Element) GetOffsetWidth() float64 {
	return js.Value(el).Get("offsetWidth").Float()
}

func (el Element) SetOffsetWidth(val float64) {
	js.Value(el).Set("offsetWidth", val)
}


func (el Element) GetOuterText() string {
	return js.Value(el).Get("outerText").String()
}

func (el Element) SetOuterText(val string) {
	js.Value(el).Set("outerText", val)
}


func (el Element) GetPopover() string {
	return js.Value(el).Get("popover").String()
}

func (el Element) SetPopover(val string) {
	js.Value(el).Set("popover", val)
}


func (el Element) GetSpellcheck() bool {
	return js.Value(el).Get("spellcheck").Bool()
}

func (el Element) SetSpellcheck(val bool) {
	js.Value(el).Set("spellcheck", val)
}


func (el Element) GetTitle() string {
	return js.Value(el).Get("title").String()
}

func (el Element) SetTitle(val string) {
	js.Value(el).Set("title", val)
}


func (el Element) GetTranslate() bool {
	return js.Value(el).Get("translate").Bool()
}

func (el Element) SetTranslate(val bool) {
	js.Value(el).Set("translate", val)
}


func (el Element) GetAttributes() js.Value {
	return js.Value(el).Get("attributes")
}

func (el Element) SetAttributes(val any) {
	js.Value(el).Set("attributes", val)
}


func (el Element) GetClassList() js.Value {
	return js.Value(el).Get("classList")
}

func (el Element) SetClassList(val any) {
	js.Value(el).Set("classList", val)
}


func (el Element) GetClassName() string {
	return js.Value(el).Get("className").String()
}

func (el Element) SetClassName(val string) {
	js.Value(el).Set("className", val)
}


func (el Element) GetClientHeight() float64 {
	return js.Value(el).Get("clientHeight").Float()
}

func (el Element) SetClientHeight(val float64) {
	js.Value(el).Set("clientHeight", val)
}


func (el Element) GetClientLeft() float64 {
	return js.Value(el).Get("clientLeft").Float()
}

func (el Element) SetClientLeft(val float64) {
	js.Value(el).Set("clientLeft", val)
}


func (el Element) GetClientTop() float64 {
	return js.Value(el).Get("clientTop").Float()
}

func (el Element) SetClientTop(val float64) {
	js.Value(el).Set("clientTop", val)
}


func (el Element) GetClientWidth() float64 {
	return js.Value(el).Get("clientWidth").Float()
}

func (el Element) SetClientWidth(val float64) {
	js.Value(el).Set("clientWidth", val)
}


func (el Element) GetId() string {
	return js.Value(el).Get("id").String()
}

func (el Element) SetId(val string) {
	js.Value(el).Set("id", val)
}


func (el Element) GetLocalName() string {
	return js.Value(el).Get("localName").String()
}

func (el Element) SetLocalName(val string) {
	js.Value(el).Set("localName", val)
}


func (el Element) GetNamespaceURI() string {
	return js.Value(el).Get("namespaceURI").String()
}

func (el Element) SetNamespaceURI(val string) {
	js.Value(el).Set("namespaceURI", val)
}


func (el Element) GetOnfullscreenchange() js.Value {
	return js.Value(el).Get("onfullscreenchange")
}

func (el Element) SetOnfullscreenchange(val any) {
	js.Value(el).Set("onfullscreenchange", val)
}


func (el Element) GetOnfullscreenerror() js.Value {
	return js.Value(el).Get("onfullscreenerror")
}

func (el Element) SetOnfullscreenerror(val any) {
	js.Value(el).Set("onfullscreenerror", val)
}


func (el Element) GetOuterHTML() string {
	return js.Value(el).Get("outerHTML").String()
}

func (el Element) SetOuterHTML(val string) {
	js.Value(el).Set("outerHTML", val)
}


func (el Element) GetOwnerDocument() js.Value {
	return js.Value(el).Get("ownerDocument")
}

func (el Element) SetOwnerDocument(val any) {
	js.Value(el).Set("ownerDocument", val)
}


func (el Element) GetPart() js.Value {
	return js.Value(el).Get("part")
}

func (el Element) SetPart(val any) {
	js.Value(el).Set("part", val)
}


func (el Element) GetPrefix() string {
	return js.Value(el).Get("prefix").String()
}

func (el Element) SetPrefix(val string) {
	js.Value(el).Set("prefix", val)
}


func (el Element) GetScrollHeight() float64 {
	return js.Value(el).Get("scrollHeight").Float()
}

func (el Element) SetScrollHeight(val float64) {
	js.Value(el).Set("scrollHeight", val)
}


func (el Element) GetScrollLeft() float64 {
	return js.Value(el).Get("scrollLeft").Float()
}

func (el Element) SetScrollLeft(val float64) {
	js.Value(el).Set("scrollLeft", val)
}


func (el Element) GetScrollTop() float64 {
	return js.Value(el).Get("scrollTop").Float()
}

func (el Element) SetScrollTop(val float64) {
	js.Value(el).Set("scrollTop", val)
}


func (el Element) GetScrollWidth() float64 {
	return js.Value(el).Get("scrollWidth").Float()
}

func (el Element) SetScrollWidth(val float64) {
	js.Value(el).Set("scrollWidth", val)
}


func (el Element) GetShadowRoot() js.Value {
	return js.Value(el).Get("shadowRoot")
}

func (el Element) SetShadowRoot(val any) {
	js.Value(el).Set("shadowRoot", val)
}


func (el Element) GetSlot() string {
	return js.Value(el).Get("slot").String()
}

func (el Element) SetSlot(val string) {
	js.Value(el).Set("slot", val)
}


func (el Element) GetTagName() string {
	return js.Value(el).Get("tagName").String()
}

func (el Element) SetTagName(val string) {
	js.Value(el).Set("tagName", val)
}


func (el Element) GetBaseURI() string {
	return js.Value(el).Get("baseURI").String()
}

func (el Element) SetBaseURI(val string) {
	js.Value(el).Set("baseURI", val)
}


func (el Element) GetChildNodes() js.Value {
	return js.Value(el).Get("childNodes")
}

func (el Element) SetChildNodes(val any) {
	js.Value(el).Set("childNodes", val)
}


func (el Element) GetFirstChild() js.Value {
	return js.Value(el).Get("firstChild")
}

func (el Element) SetFirstChild(val any) {
	js.Value(el).Set("firstChild", val)
}


func (el Element) GetIsConnected() bool {
	return js.Value(el).Get("isConnected").Bool()
}

func (el Element) SetIsConnected(val bool) {
	js.Value(el).Set("isConnected", val)
}


func (el Element) GetLastChild() js.Value {
	return js.Value(el).Get("lastChild")
}

func (el Element) SetLastChild(val any) {
	js.Value(el).Set("lastChild", val)
}


func (el Element) GetNextSibling() js.Value {
	return js.Value(el).Get("nextSibling")
}

func (el Element) SetNextSibling(val any) {
	js.Value(el).Set("nextSibling", val)
}


func (el Element) GetNodeName() string {
	return js.Value(el).Get("nodeName").String()
}

func (el Element) SetNodeName(val string) {
	js.Value(el).Set("nodeName", val)
}


func (el Element) GetNodeType() float64 {
	return js.Value(el).Get("nodeType").Float()
}

func (el Element) SetNodeType(val float64) {
	js.Value(el).Set("nodeType", val)
}


func (el Element) GetNodeValue() string {
	return js.Value(el).Get("nodeValue").String()
}

func (el Element) SetNodeValue(val string) {
	js.Value(el).Set("nodeValue", val)
}


func (el Element) GetParentElement() js.Value {
	return js.Value(el).Get("parentElement")
}

func (el Element) SetParentElement(val any) {
	js.Value(el).Set("parentElement", val)
}


func (el Element) GetParentNode() js.Value {
	return js.Value(el).Get("parentNode")
}

func (el Element) SetParentNode(val any) {
	js.Value(el).Set("parentNode", val)
}


func (el Element) GetPreviousSibling() js.Value {
	return js.Value(el).Get("previousSibling")
}

func (el Element) SetPreviousSibling(val any) {
	js.Value(el).Set("previousSibling", val)
}


func (el Element) GetTextContent() string {
	return js.Value(el).Get("textContent").String()
}

func (el Element) SetTextContent(val string) {
	js.Value(el).Set("textContent", val)
}


func (el Element) GetELEMENT_NODE() js.Value {
	return js.Value(el).Get("ELEMENT_NODE")
}

func (el Element) SetELEMENT_NODE(val any) {
	js.Value(el).Set("ELEMENT_NODE", val)
}


func (el Element) GetATTRIBUTE_NODE() js.Value {
	return js.Value(el).Get("ATTRIBUTE_NODE")
}

func (el Element) SetATTRIBUTE_NODE(val any) {
	js.Value(el).Set("ATTRIBUTE_NODE", val)
}


func (el Element) GetTEXT_NODE() js.Value {
	return js.Value(el).Get("TEXT_NODE")
}

func (el Element) SetTEXT_NODE(val any) {
	js.Value(el).Set("TEXT_NODE", val)
}


func (el Element) GetCDATA_SECTION_NODE() js.Value {
	return js.Value(el).Get("CDATA_SECTION_NODE")
}

func (el Element) SetCDATA_SECTION_NODE(val any) {
	js.Value(el).Set("CDATA_SECTION_NODE", val)
}


func (el Element) GetENTITY_REFERENCE_NODE() js.Value {
	return js.Value(el).Get("ENTITY_REFERENCE_NODE")
}

func (el Element) SetENTITY_REFERENCE_NODE(val any) {
	js.Value(el).Set("ENTITY_REFERENCE_NODE", val)
}


func (el Element) GetENTITY_NODE() js.Value {
	return js.Value(el).Get("ENTITY_NODE")
}

func (el Element) SetENTITY_NODE(val any) {
	js.Value(el).Set("ENTITY_NODE", val)
}


func (el Element) GetPROCESSING_INSTRUCTION_NODE() js.Value {
	return js.Value(el).Get("PROCESSING_INSTRUCTION_NODE")
}

func (el Element) SetPROCESSING_INSTRUCTION_NODE(val any) {
	js.Value(el).Set("PROCESSING_INSTRUCTION_NODE", val)
}


func (el Element) GetCOMMENT_NODE() js.Value {
	return js.Value(el).Get("COMMENT_NODE")
}

func (el Element) SetCOMMENT_NODE(val any) {
	js.Value(el).Set("COMMENT_NODE", val)
}


func (el Element) GetDOCUMENT_NODE() js.Value {
	return js.Value(el).Get("DOCUMENT_NODE")
}

func (el Element) SetDOCUMENT_NODE(val any) {
	js.Value(el).Set("DOCUMENT_NODE", val)
}


func (el Element) GetDOCUMENT_TYPE_NODE() js.Value {
	return js.Value(el).Get("DOCUMENT_TYPE_NODE")
}

func (el Element) SetDOCUMENT_TYPE_NODE(val any) {
	js.Value(el).Set("DOCUMENT_TYPE_NODE", val)
}


func (el Element) GetDOCUMENT_FRAGMENT_NODE() js.Value {
	return js.Value(el).Get("DOCUMENT_FRAGMENT_NODE")
}

func (el Element) SetDOCUMENT_FRAGMENT_NODE(val any) {
	js.Value(el).Set("DOCUMENT_FRAGMENT_NODE", val)
}


func (el Element) GetNOTATION_NODE() js.Value {
	return js.Value(el).Get("NOTATION_NODE")
}

func (el Element) SetNOTATION_NODE(val any) {
	js.Value(el).Set("NOTATION_NODE", val)
}


func (el Element) GetDOCUMENT_POSITION_DISCONNECTED() js.Value {
	return js.Value(el).Get("DOCUMENT_POSITION_DISCONNECTED")
}

func (el Element) SetDOCUMENT_POSITION_DISCONNECTED(val any) {
	js.Value(el).Set("DOCUMENT_POSITION_DISCONNECTED", val)
}


func (el Element) GetDOCUMENT_POSITION_PRECEDING() js.Value {
	return js.Value(el).Get("DOCUMENT_POSITION_PRECEDING")
}

func (el Element) SetDOCUMENT_POSITION_PRECEDING(val any) {
	js.Value(el).Set("DOCUMENT_POSITION_PRECEDING", val)
}


func (el Element) GetDOCUMENT_POSITION_FOLLOWING() js.Value {
	return js.Value(el).Get("DOCUMENT_POSITION_FOLLOWING")
}

func (el Element) SetDOCUMENT_POSITION_FOLLOWING(val any) {
	js.Value(el).Set("DOCUMENT_POSITION_FOLLOWING", val)
}


func (el Element) GetDOCUMENT_POSITION_CONTAINS() js.Value {
	return js.Value(el).Get("DOCUMENT_POSITION_CONTAINS")
}

func (el Element) SetDOCUMENT_POSITION_CONTAINS(val any) {
	js.Value(el).Set("DOCUMENT_POSITION_CONTAINS", val)
}


func (el Element) GetDOCUMENT_POSITION_CONTAINED_BY() js.Value {
	return js.Value(el).Get("DOCUMENT_POSITION_CONTAINED_BY")
}

func (el Element) SetDOCUMENT_POSITION_CONTAINED_BY(val any) {
	js.Value(el).Set("DOCUMENT_POSITION_CONTAINED_BY", val)
}


func (el Element) GetDOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC() js.Value {
	return js.Value(el).Get("DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC")
}

func (el Element) SetDOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC(val any) {
	js.Value(el).Set("DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC", val)
}


func (el Element) GetAriaAtomic() string {
	return js.Value(el).Get("ariaAtomic").String()
}

func (el Element) SetAriaAtomic(val string) {
	js.Value(el).Set("ariaAtomic", val)
}


func (el Element) GetAriaAutoComplete() string {
	return js.Value(el).Get("ariaAutoComplete").String()
}

func (el Element) SetAriaAutoComplete(val string) {
	js.Value(el).Set("ariaAutoComplete", val)
}


func (el Element) GetAriaBusy() string {
	return js.Value(el).Get("ariaBusy").String()
}

func (el Element) SetAriaBusy(val string) {
	js.Value(el).Set("ariaBusy", val)
}


func (el Element) GetAriaChecked() string {
	return js.Value(el).Get("ariaChecked").String()
}

func (el Element) SetAriaChecked(val string) {
	js.Value(el).Set("ariaChecked", val)
}


func (el Element) GetAriaColCount() string {
	return js.Value(el).Get("ariaColCount").String()
}

func (el Element) SetAriaColCount(val string) {
	js.Value(el).Set("ariaColCount", val)
}


func (el Element) GetAriaColIndex() string {
	return js.Value(el).Get("ariaColIndex").String()
}

func (el Element) SetAriaColIndex(val string) {
	js.Value(el).Set("ariaColIndex", val)
}


func (el Element) GetAriaColSpan() string {
	return js.Value(el).Get("ariaColSpan").String()
}

func (el Element) SetAriaColSpan(val string) {
	js.Value(el).Set("ariaColSpan", val)
}


func (el Element) GetAriaCurrent() string {
	return js.Value(el).Get("ariaCurrent").String()
}

func (el Element) SetAriaCurrent(val string) {
	js.Value(el).Set("ariaCurrent", val)
}


func (el Element) GetAriaDisabled() string {
	return js.Value(el).Get("ariaDisabled").String()
}

func (el Element) SetAriaDisabled(val string) {
	js.Value(el).Set("ariaDisabled", val)
}


func (el Element) GetAriaExpanded() string {
	return js.Value(el).Get("ariaExpanded").String()
}

func (el Element) SetAriaExpanded(val string) {
	js.Value(el).Set("ariaExpanded", val)
}


func (el Element) GetAriaHasPopup() string {
	return js.Value(el).Get("ariaHasPopup").String()
}

func (el Element) SetAriaHasPopup(val string) {
	js.Value(el).Set("ariaHasPopup", val)
}


func (el Element) GetAriaHidden() string {
	return js.Value(el).Get("ariaHidden").String()
}

func (el Element) SetAriaHidden(val string) {
	js.Value(el).Set("ariaHidden", val)
}


func (el Element) GetAriaInvalid() string {
	return js.Value(el).Get("ariaInvalid").String()
}

func (el Element) SetAriaInvalid(val string) {
	js.Value(el).Set("ariaInvalid", val)
}


func (el Element) GetAriaKeyShortcuts() string {
	return js.Value(el).Get("ariaKeyShortcuts").String()
}

func (el Element) SetAriaKeyShortcuts(val string) {
	js.Value(el).Set("ariaKeyShortcuts", val)
}


func (el Element) GetAriaLabel() string {
	return js.Value(el).Get("ariaLabel").String()
}

func (el Element) SetAriaLabel(val string) {
	js.Value(el).Set("ariaLabel", val)
}


func (el Element) GetAriaLevel() string {
	return js.Value(el).Get("ariaLevel").String()
}

func (el Element) SetAriaLevel(val string) {
	js.Value(el).Set("ariaLevel", val)
}


func (el Element) GetAriaLive() string {
	return js.Value(el).Get("ariaLive").String()
}

func (el Element) SetAriaLive(val string) {
	js.Value(el).Set("ariaLive", val)
}


func (el Element) GetAriaModal() string {
	return js.Value(el).Get("ariaModal").String()
}

func (el Element) SetAriaModal(val string) {
	js.Value(el).Set("ariaModal", val)
}


func (el Element) GetAriaMultiLine() string {
	return js.Value(el).Get("ariaMultiLine").String()
}

func (el Element) SetAriaMultiLine(val string) {
	js.Value(el).Set("ariaMultiLine", val)
}


func (el Element) GetAriaMultiSelectable() string {
	return js.Value(el).Get("ariaMultiSelectable").String()
}

func (el Element) SetAriaMultiSelectable(val string) {
	js.Value(el).Set("ariaMultiSelectable", val)
}


func (el Element) GetAriaOrientation() string {
	return js.Value(el).Get("ariaOrientation").String()
}

func (el Element) SetAriaOrientation(val string) {
	js.Value(el).Set("ariaOrientation", val)
}


func (el Element) GetAriaPlaceholder() string {
	return js.Value(el).Get("ariaPlaceholder").String()
}

func (el Element) SetAriaPlaceholder(val string) {
	js.Value(el).Set("ariaPlaceholder", val)
}


func (el Element) GetAriaPosInSet() string {
	return js.Value(el).Get("ariaPosInSet").String()
}

func (el Element) SetAriaPosInSet(val string) {
	js.Value(el).Set("ariaPosInSet", val)
}


func (el Element) GetAriaPressed() string {
	return js.Value(el).Get("ariaPressed").String()
}

func (el Element) SetAriaPressed(val string) {
	js.Value(el).Set("ariaPressed", val)
}


func (el Element) GetAriaReadOnly() string {
	return js.Value(el).Get("ariaReadOnly").String()
}

func (el Element) SetAriaReadOnly(val string) {
	js.Value(el).Set("ariaReadOnly", val)
}


func (el Element) GetAriaRequired() string {
	return js.Value(el).Get("ariaRequired").String()
}

func (el Element) SetAriaRequired(val string) {
	js.Value(el).Set("ariaRequired", val)
}


func (el Element) GetAriaRoleDescription() string {
	return js.Value(el).Get("ariaRoleDescription").String()
}

func (el Element) SetAriaRoleDescription(val string) {
	js.Value(el).Set("ariaRoleDescription", val)
}


func (el Element) GetAriaRowCount() string {
	return js.Value(el).Get("ariaRowCount").String()
}

func (el Element) SetAriaRowCount(val string) {
	js.Value(el).Set("ariaRowCount", val)
}


func (el Element) GetAriaRowIndex() string {
	return js.Value(el).Get("ariaRowIndex").String()
}

func (el Element) SetAriaRowIndex(val string) {
	js.Value(el).Set("ariaRowIndex", val)
}


func (el Element) GetAriaRowSpan() string {
	return js.Value(el).Get("ariaRowSpan").String()
}

func (el Element) SetAriaRowSpan(val string) {
	js.Value(el).Set("ariaRowSpan", val)
}


func (el Element) GetAriaSelected() string {
	return js.Value(el).Get("ariaSelected").String()
}

func (el Element) SetAriaSelected(val string) {
	js.Value(el).Set("ariaSelected", val)
}


func (el Element) GetAriaSetSize() string {
	return js.Value(el).Get("ariaSetSize").String()
}

func (el Element) SetAriaSetSize(val string) {
	js.Value(el).Set("ariaSetSize", val)
}


func (el Element) GetAriaSort() string {
	return js.Value(el).Get("ariaSort").String()
}

func (el Element) SetAriaSort(val string) {
	js.Value(el).Set("ariaSort", val)
}


func (el Element) GetAriaValueMax() string {
	return js.Value(el).Get("ariaValueMax").String()
}

func (el Element) SetAriaValueMax(val string) {
	js.Value(el).Set("ariaValueMax", val)
}


func (el Element) GetAriaValueMin() string {
	return js.Value(el).Get("ariaValueMin").String()
}

func (el Element) SetAriaValueMin(val string) {
	js.Value(el).Set("ariaValueMin", val)
}


func (el Element) GetAriaValueNow() string {
	return js.Value(el).Get("ariaValueNow").String()
}

func (el Element) SetAriaValueNow(val string) {
	js.Value(el).Set("ariaValueNow", val)
}


func (el Element) GetAriaValueText() string {
	return js.Value(el).Get("ariaValueText").String()
}

func (el Element) SetAriaValueText(val string) {
	js.Value(el).Set("ariaValueText", val)
}


func (el Element) GetRole() string {
	return js.Value(el).Get("role").String()
}

func (el Element) SetRole(val string) {
	js.Value(el).Set("role", val)
}


func (el Element) GetInnerHTML() string {
	return js.Value(el).Get("innerHTML").String()
}

func (el Element) SetInnerHTML(val string) {
	js.Value(el).Set("innerHTML", val)
}


func (el Element) GetNextElementSibling() js.Value {
	return js.Value(el).Get("nextElementSibling")
}

func (el Element) SetNextElementSibling(val any) {
	js.Value(el).Set("nextElementSibling", val)
}


func (el Element) GetPreviousElementSibling() js.Value {
	return js.Value(el).Get("previousElementSibling")
}

func (el Element) SetPreviousElementSibling(val any) {
	js.Value(el).Set("previousElementSibling", val)
}


func (el Element) GetChildElementCount() float64 {
	return js.Value(el).Get("childElementCount").Float()
}

func (el Element) SetChildElementCount(val float64) {
	js.Value(el).Set("childElementCount", val)
}


func (el Element) GetChildren() js.Value {
	return js.Value(el).Get("children")
}

func (el Element) SetChildren(val any) {
	js.Value(el).Set("children", val)
}


func (el Element) GetFirstElementChild() js.Value {
	return js.Value(el).Get("firstElementChild")
}

func (el Element) SetFirstElementChild(val any) {
	js.Value(el).Set("firstElementChild", val)
}


func (el Element) GetLastElementChild() js.Value {
	return js.Value(el).Get("lastElementChild")
}

func (el Element) SetLastElementChild(val any) {
	js.Value(el).Set("lastElementChild", val)
}


func (el Element) GetAssignedSlot() js.Value {
	return js.Value(el).Get("assignedSlot")
}

func (el Element) SetAssignedSlot(val any) {
	js.Value(el).Set("assignedSlot", val)
}


func (el Element) GetAttributeStyleMap() js.Value {
	return js.Value(el).Get("attributeStyleMap")
}

func (el Element) SetAttributeStyleMap(val any) {
	js.Value(el).Set("attributeStyleMap", val)
}


func (el Element) GetStyle() js.Value {
	return js.Value(el).Get("style")
}

func (el Element) SetStyle(val any) {
	js.Value(el).Set("style", val)
}


func (el Element) GetContentEditable() string {
	return js.Value(el).Get("contentEditable").String()
}

func (el Element) SetContentEditable(val string) {
	js.Value(el).Set("contentEditable", val)
}


func (el Element) GetEnterKeyHint() string {
	return js.Value(el).Get("enterKeyHint").String()
}

func (el Element) SetEnterKeyHint(val string) {
	js.Value(el).Set("enterKeyHint", val)
}


func (el Element) GetInputMode() string {
	return js.Value(el).Get("inputMode").String()
}

func (el Element) SetInputMode(val string) {
	js.Value(el).Set("inputMode", val)
}


func (el Element) GetIsContentEditable() bool {
	return js.Value(el).Get("isContentEditable").Bool()
}

func (el Element) SetIsContentEditable(val bool) {
	js.Value(el).Set("isContentEditable", val)
}


func (el Element) GetOnabort() js.Value {
	return js.Value(el).Get("onabort")
}

func (el Element) SetOnabort(val any) {
	js.Value(el).Set("onabort", val)
}


func (el Element) GetOnanimationcancel() js.Value {
	return js.Value(el).Get("onanimationcancel")
}

func (el Element) SetOnanimationcancel(val any) {
	js.Value(el).Set("onanimationcancel", val)
}


func (el Element) GetOnanimationend() js.Value {
	return js.Value(el).Get("onanimationend")
}

func (el Element) SetOnanimationend(val any) {
	js.Value(el).Set("onanimationend", val)
}


func (el Element) GetOnanimationiteration() js.Value {
	return js.Value(el).Get("onanimationiteration")
}

func (el Element) SetOnanimationiteration(val any) {
	js.Value(el).Set("onanimationiteration", val)
}


func (el Element) GetOnanimationstart() js.Value {
	return js.Value(el).Get("onanimationstart")
}

func (el Element) SetOnanimationstart(val any) {
	js.Value(el).Set("onanimationstart", val)
}


func (el Element) GetOnauxclick() js.Value {
	return js.Value(el).Get("onauxclick")
}

func (el Element) SetOnauxclick(val any) {
	js.Value(el).Set("onauxclick", val)
}


func (el Element) GetOnbeforeinput() js.Value {
	return js.Value(el).Get("onbeforeinput")
}

func (el Element) SetOnbeforeinput(val any) {
	js.Value(el).Set("onbeforeinput", val)
}


func (el Element) GetOnblur() js.Value {
	return js.Value(el).Get("onblur")
}

func (el Element) SetOnblur(val any) {
	js.Value(el).Set("onblur", val)
}


func (el Element) GetOncancel() js.Value {
	return js.Value(el).Get("oncancel")
}

func (el Element) SetOncancel(val any) {
	js.Value(el).Set("oncancel", val)
}


func (el Element) GetOncanplay() js.Value {
	return js.Value(el).Get("oncanplay")
}

func (el Element) SetOncanplay(val any) {
	js.Value(el).Set("oncanplay", val)
}


func (el Element) GetOncanplaythrough() js.Value {
	return js.Value(el).Get("oncanplaythrough")
}

func (el Element) SetOncanplaythrough(val any) {
	js.Value(el).Set("oncanplaythrough", val)
}


func (el Element) GetOnchange() js.Value {
	return js.Value(el).Get("onchange")
}

func (el Element) SetOnchange(val any) {
	js.Value(el).Set("onchange", val)
}


func (el Element) GetOnclick() js.Value {
	return js.Value(el).Get("onclick")
}

func (el Element) SetOnclick(val any) {
	js.Value(el).Set("onclick", val)
}


func (el Element) GetOnclose() js.Value {
	return js.Value(el).Get("onclose")
}

func (el Element) SetOnclose(val any) {
	js.Value(el).Set("onclose", val)
}


func (el Element) GetOncontextmenu() js.Value {
	return js.Value(el).Get("oncontextmenu")
}

func (el Element) SetOncontextmenu(val any) {
	js.Value(el).Set("oncontextmenu", val)
}


func (el Element) GetOncopy() js.Value {
	return js.Value(el).Get("oncopy")
}

func (el Element) SetOncopy(val any) {
	js.Value(el).Set("oncopy", val)
}


func (el Element) GetOncuechange() js.Value {
	return js.Value(el).Get("oncuechange")
}

func (el Element) SetOncuechange(val any) {
	js.Value(el).Set("oncuechange", val)
}


func (el Element) GetOncut() js.Value {
	return js.Value(el).Get("oncut")
}

func (el Element) SetOncut(val any) {
	js.Value(el).Set("oncut", val)
}


func (el Element) GetOndblclick() js.Value {
	return js.Value(el).Get("ondblclick")
}

func (el Element) SetOndblclick(val any) {
	js.Value(el).Set("ondblclick", val)
}


func (el Element) GetOndrag() js.Value {
	return js.Value(el).Get("ondrag")
}

func (el Element) SetOndrag(val any) {
	js.Value(el).Set("ondrag", val)
}


func (el Element) GetOndragend() js.Value {
	return js.Value(el).Get("ondragend")
}

func (el Element) SetOndragend(val any) {
	js.Value(el).Set("ondragend", val)
}


func (el Element) GetOndragenter() js.Value {
	return js.Value(el).Get("ondragenter")
}

func (el Element) SetOndragenter(val any) {
	js.Value(el).Set("ondragenter", val)
}


func (el Element) GetOndragleave() js.Value {
	return js.Value(el).Get("ondragleave")
}

func (el Element) SetOndragleave(val any) {
	js.Value(el).Set("ondragleave", val)
}


func (el Element) GetOndragover() js.Value {
	return js.Value(el).Get("ondragover")
}

func (el Element) SetOndragover(val any) {
	js.Value(el).Set("ondragover", val)
}


func (el Element) GetOndragstart() js.Value {
	return js.Value(el).Get("ondragstart")
}

func (el Element) SetOndragstart(val any) {
	js.Value(el).Set("ondragstart", val)
}


func (el Element) GetOndrop() js.Value {
	return js.Value(el).Get("ondrop")
}

func (el Element) SetOndrop(val any) {
	js.Value(el).Set("ondrop", val)
}


func (el Element) GetOndurationchange() js.Value {
	return js.Value(el).Get("ondurationchange")
}

func (el Element) SetOndurationchange(val any) {
	js.Value(el).Set("ondurationchange", val)
}


func (el Element) GetOnemptied() js.Value {
	return js.Value(el).Get("onemptied")
}

func (el Element) SetOnemptied(val any) {
	js.Value(el).Set("onemptied", val)
}


func (el Element) GetOnended() js.Value {
	return js.Value(el).Get("onended")
}

func (el Element) SetOnended(val any) {
	js.Value(el).Set("onended", val)
}


func (el Element) GetOnerror() js.Value {
	return js.Value(el).Get("onerror")
}

func (el Element) SetOnerror(val any) {
	js.Value(el).Set("onerror", val)
}


func (el Element) GetOnfocus() js.Value {
	return js.Value(el).Get("onfocus")
}

func (el Element) SetOnfocus(val any) {
	js.Value(el).Set("onfocus", val)
}


func (el Element) GetOnformdata() js.Value {
	return js.Value(el).Get("onformdata")
}

func (el Element) SetOnformdata(val any) {
	js.Value(el).Set("onformdata", val)
}


func (el Element) GetOngotpointercapture() js.Value {
	return js.Value(el).Get("ongotpointercapture")
}

func (el Element) SetOngotpointercapture(val any) {
	js.Value(el).Set("ongotpointercapture", val)
}


func (el Element) GetOninput() js.Value {
	return js.Value(el).Get("oninput")
}

func (el Element) SetOninput(val any) {
	js.Value(el).Set("oninput", val)
}


func (el Element) GetOninvalid() js.Value {
	return js.Value(el).Get("oninvalid")
}

func (el Element) SetOninvalid(val any) {
	js.Value(el).Set("oninvalid", val)
}


func (el Element) GetOnkeydown() js.Value {
	return js.Value(el).Get("onkeydown")
}

func (el Element) SetOnkeydown(val any) {
	js.Value(el).Set("onkeydown", val)
}


func (el Element) GetOnkeypress() js.Value {
	return js.Value(el).Get("onkeypress")
}

func (el Element) SetOnkeypress(val any) {
	js.Value(el).Set("onkeypress", val)
}


func (el Element) GetOnkeyup() js.Value {
	return js.Value(el).Get("onkeyup")
}

func (el Element) SetOnkeyup(val any) {
	js.Value(el).Set("onkeyup", val)
}


func (el Element) GetOnload() js.Value {
	return js.Value(el).Get("onload")
}

func (el Element) SetOnload(val any) {
	js.Value(el).Set("onload", val)
}


func (el Element) GetOnloadeddata() js.Value {
	return js.Value(el).Get("onloadeddata")
}

func (el Element) SetOnloadeddata(val any) {
	js.Value(el).Set("onloadeddata", val)
}


func (el Element) GetOnloadedmetadata() js.Value {
	return js.Value(el).Get("onloadedmetadata")
}

func (el Element) SetOnloadedmetadata(val any) {
	js.Value(el).Set("onloadedmetadata", val)
}


func (el Element) GetOnloadstart() js.Value {
	return js.Value(el).Get("onloadstart")
}

func (el Element) SetOnloadstart(val any) {
	js.Value(el).Set("onloadstart", val)
}


func (el Element) GetOnlostpointercapture() js.Value {
	return js.Value(el).Get("onlostpointercapture")
}

func (el Element) SetOnlostpointercapture(val any) {
	js.Value(el).Set("onlostpointercapture", val)
}


func (el Element) GetOnmousedown() js.Value {
	return js.Value(el).Get("onmousedown")
}

func (el Element) SetOnmousedown(val any) {
	js.Value(el).Set("onmousedown", val)
}


func (el Element) GetOnmouseenter() js.Value {
	return js.Value(el).Get("onmouseenter")
}

func (el Element) SetOnmouseenter(val any) {
	js.Value(el).Set("onmouseenter", val)
}


func (el Element) GetOnmouseleave() js.Value {
	return js.Value(el).Get("onmouseleave")
}

func (el Element) SetOnmouseleave(val any) {
	js.Value(el).Set("onmouseleave", val)
}


func (el Element) GetOnmousemove() js.Value {
	return js.Value(el).Get("onmousemove")
}

func (el Element) SetOnmousemove(val any) {
	js.Value(el).Set("onmousemove", val)
}


func (el Element) GetOnmouseout() js.Value {
	return js.Value(el).Get("onmouseout")
}

func (el Element) SetOnmouseout(val any) {
	js.Value(el).Set("onmouseout", val)
}


func (el Element) GetOnmouseover() js.Value {
	return js.Value(el).Get("onmouseover")
}

func (el Element) SetOnmouseover(val any) {
	js.Value(el).Set("onmouseover", val)
}


func (el Element) GetOnmouseup() js.Value {
	return js.Value(el).Get("onmouseup")
}

func (el Element) SetOnmouseup(val any) {
	js.Value(el).Set("onmouseup", val)
}


func (el Element) GetOnpaste() js.Value {
	return js.Value(el).Get("onpaste")
}

func (el Element) SetOnpaste(val any) {
	js.Value(el).Set("onpaste", val)
}


func (el Element) GetOnpause() js.Value {
	return js.Value(el).Get("onpause")
}

func (el Element) SetOnpause(val any) {
	js.Value(el).Set("onpause", val)
}


func (el Element) GetOnplay() js.Value {
	return js.Value(el).Get("onplay")
}

func (el Element) SetOnplay(val any) {
	js.Value(el).Set("onplay", val)
}


func (el Element) GetOnplaying() js.Value {
	return js.Value(el).Get("onplaying")
}

func (el Element) SetOnplaying(val any) {
	js.Value(el).Set("onplaying", val)
}


func (el Element) GetOnpointercancel() js.Value {
	return js.Value(el).Get("onpointercancel")
}

func (el Element) SetOnpointercancel(val any) {
	js.Value(el).Set("onpointercancel", val)
}


func (el Element) GetOnpointerdown() js.Value {
	return js.Value(el).Get("onpointerdown")
}

func (el Element) SetOnpointerdown(val any) {
	js.Value(el).Set("onpointerdown", val)
}


func (el Element) GetOnpointerenter() js.Value {
	return js.Value(el).Get("onpointerenter")
}

func (el Element) SetOnpointerenter(val any) {
	js.Value(el).Set("onpointerenter", val)
}


func (el Element) GetOnpointerleave() js.Value {
	return js.Value(el).Get("onpointerleave")
}

func (el Element) SetOnpointerleave(val any) {
	js.Value(el).Set("onpointerleave", val)
}


func (el Element) GetOnpointermove() js.Value {
	return js.Value(el).Get("onpointermove")
}

func (el Element) SetOnpointermove(val any) {
	js.Value(el).Set("onpointermove", val)
}


func (el Element) GetOnpointerout() js.Value {
	return js.Value(el).Get("onpointerout")
}

func (el Element) SetOnpointerout(val any) {
	js.Value(el).Set("onpointerout", val)
}


func (el Element) GetOnpointerover() js.Value {
	return js.Value(el).Get("onpointerover")
}

func (el Element) SetOnpointerover(val any) {
	js.Value(el).Set("onpointerover", val)
}


func (el Element) GetOnpointerup() js.Value {
	return js.Value(el).Get("onpointerup")
}

func (el Element) SetOnpointerup(val any) {
	js.Value(el).Set("onpointerup", val)
}


func (el Element) GetOnprogress() js.Value {
	return js.Value(el).Get("onprogress")
}

func (el Element) SetOnprogress(val any) {
	js.Value(el).Set("onprogress", val)
}


func (el Element) GetOnratechange() js.Value {
	return js.Value(el).Get("onratechange")
}

func (el Element) SetOnratechange(val any) {
	js.Value(el).Set("onratechange", val)
}


func (el Element) GetOnreset() js.Value {
	return js.Value(el).Get("onreset")
}

func (el Element) SetOnreset(val any) {
	js.Value(el).Set("onreset", val)
}


func (el Element) GetOnresize() js.Value {
	return js.Value(el).Get("onresize")
}

func (el Element) SetOnresize(val any) {
	js.Value(el).Set("onresize", val)
}


func (el Element) GetOnscroll() js.Value {
	return js.Value(el).Get("onscroll")
}

func (el Element) SetOnscroll(val any) {
	js.Value(el).Set("onscroll", val)
}


func (el Element) GetOnscrollend() js.Value {
	return js.Value(el).Get("onscrollend")
}

func (el Element) SetOnscrollend(val any) {
	js.Value(el).Set("onscrollend", val)
}


func (el Element) GetOnsecuritypolicyviolation() js.Value {
	return js.Value(el).Get("onsecuritypolicyviolation")
}

func (el Element) SetOnsecuritypolicyviolation(val any) {
	js.Value(el).Set("onsecuritypolicyviolation", val)
}


func (el Element) GetOnseeked() js.Value {
	return js.Value(el).Get("onseeked")
}

func (el Element) SetOnseeked(val any) {
	js.Value(el).Set("onseeked", val)
}


func (el Element) GetOnseeking() js.Value {
	return js.Value(el).Get("onseeking")
}

func (el Element) SetOnseeking(val any) {
	js.Value(el).Set("onseeking", val)
}


func (el Element) GetOnselect() js.Value {
	return js.Value(el).Get("onselect")
}

func (el Element) SetOnselect(val any) {
	js.Value(el).Set("onselect", val)
}


func (el Element) GetOnselectionchange() js.Value {
	return js.Value(el).Get("onselectionchange")
}

func (el Element) SetOnselectionchange(val any) {
	js.Value(el).Set("onselectionchange", val)
}


func (el Element) GetOnselectstart() js.Value {
	return js.Value(el).Get("onselectstart")
}

func (el Element) SetOnselectstart(val any) {
	js.Value(el).Set("onselectstart", val)
}


func (el Element) GetOnslotchange() js.Value {
	return js.Value(el).Get("onslotchange")
}

func (el Element) SetOnslotchange(val any) {
	js.Value(el).Set("onslotchange", val)
}


func (el Element) GetOnstalled() js.Value {
	return js.Value(el).Get("onstalled")
}

func (el Element) SetOnstalled(val any) {
	js.Value(el).Set("onstalled", val)
}


func (el Element) GetOnsubmit() js.Value {
	return js.Value(el).Get("onsubmit")
}

func (el Element) SetOnsubmit(val any) {
	js.Value(el).Set("onsubmit", val)
}


func (el Element) GetOnsuspend() js.Value {
	return js.Value(el).Get("onsuspend")
}

func (el Element) SetOnsuspend(val any) {
	js.Value(el).Set("onsuspend", val)
}


func (el Element) GetOntimeupdate() js.Value {
	return js.Value(el).Get("ontimeupdate")
}

func (el Element) SetOntimeupdate(val any) {
	js.Value(el).Set("ontimeupdate", val)
}


func (el Element) GetOntoggle() js.Value {
	return js.Value(el).Get("ontoggle")
}

func (el Element) SetOntoggle(val any) {
	js.Value(el).Set("ontoggle", val)
}


func (el Element) GetOntouchcancel() js.Value {
	return js.Value(el).Get("ontouchcancel")
}

func (el Element) SetOntouchcancel(val any) {
	js.Value(el).Set("ontouchcancel", val)
}


func (el Element) GetOntouchend() js.Value {
	return js.Value(el).Get("ontouchend")
}

func (el Element) SetOntouchend(val any) {
	js.Value(el).Set("ontouchend", val)
}


func (el Element) GetOntouchmove() js.Value {
	return js.Value(el).Get("ontouchmove")
}

func (el Element) SetOntouchmove(val any) {
	js.Value(el).Set("ontouchmove", val)
}


func (el Element) GetOntouchstart() js.Value {
	return js.Value(el).Get("ontouchstart")
}

func (el Element) SetOntouchstart(val any) {
	js.Value(el).Set("ontouchstart", val)
}


func (el Element) GetOntransitioncancel() js.Value {
	return js.Value(el).Get("ontransitioncancel")
}

func (el Element) SetOntransitioncancel(val any) {
	js.Value(el).Set("ontransitioncancel", val)
}


func (el Element) GetOntransitionend() js.Value {
	return js.Value(el).Get("ontransitionend")
}

func (el Element) SetOntransitionend(val any) {
	js.Value(el).Set("ontransitionend", val)
}


func (el Element) GetOntransitionrun() js.Value {
	return js.Value(el).Get("ontransitionrun")
}

func (el Element) SetOntransitionrun(val any) {
	js.Value(el).Set("ontransitionrun", val)
}


func (el Element) GetOntransitionstart() js.Value {
	return js.Value(el).Get("ontransitionstart")
}

func (el Element) SetOntransitionstart(val any) {
	js.Value(el).Set("ontransitionstart", val)
}


func (el Element) GetOnvolumechange() js.Value {
	return js.Value(el).Get("onvolumechange")
}

func (el Element) SetOnvolumechange(val any) {
	js.Value(el).Set("onvolumechange", val)
}


func (el Element) GetOnwaiting() js.Value {
	return js.Value(el).Get("onwaiting")
}

func (el Element) SetOnwaiting(val any) {
	js.Value(el).Set("onwaiting", val)
}


func (el Element) GetOnwebkitanimationend() js.Value {
	return js.Value(el).Get("onwebkitanimationend")
}

func (el Element) SetOnwebkitanimationend(val any) {
	js.Value(el).Set("onwebkitanimationend", val)
}


func (el Element) GetOnwebkitanimationiteration() js.Value {
	return js.Value(el).Get("onwebkitanimationiteration")
}

func (el Element) SetOnwebkitanimationiteration(val any) {
	js.Value(el).Set("onwebkitanimationiteration", val)
}


func (el Element) GetOnwebkitanimationstart() js.Value {
	return js.Value(el).Get("onwebkitanimationstart")
}

func (el Element) SetOnwebkitanimationstart(val any) {
	js.Value(el).Set("onwebkitanimationstart", val)
}


func (el Element) GetOnwebkittransitionend() js.Value {
	return js.Value(el).Get("onwebkittransitionend")
}

func (el Element) SetOnwebkittransitionend(val any) {
	js.Value(el).Set("onwebkittransitionend", val)
}


func (el Element) GetOnwheel() js.Value {
	return js.Value(el).Get("onwheel")
}

func (el Element) SetOnwheel(val any) {
	js.Value(el).Set("onwheel", val)
}


func (el Element) GetAutofocus() bool {
	return js.Value(el).Get("autofocus").Bool()
}

func (el Element) SetAutofocus(val bool) {
	js.Value(el).Set("autofocus", val)
}


func (el Element) GetDataset() js.Value {
	return js.Value(el).Get("dataset")
}

func (el Element) SetDataset(val any) {
	js.Value(el).Set("dataset", val)
}


func (el Element) GetNonce() string {
	return js.Value(el).Get("nonce").String()
}

func (el Element) SetNonce(val string) {
	js.Value(el).Set("nonce", val)
}


func (el Element) GetTabIndex() float64 {
	return js.Value(el).Get("tabIndex").Float()
}

func (el Element) SetTabIndex(val float64) {
	js.Value(el).Set("tabIndex", val)
}



func (el Element) AttachInternals() js.Value {
	return js.Value(el).Call("attachInternals")
}


func (el Element) Click() js.Value {
	return js.Value(el).Call("click")
}


func (el Element) HidePopover() js.Value {
	return js.Value(el).Call("hidePopover")
}


func (el Element) ShowPopover() js.Value {
	return js.Value(el).Call("showPopover")
}


func (el Element) TogglePopover0() js.Value {
	return js.Value(el).Call("togglePopover")
}


func (el Element) TogglePopover1(force bool) js.Value {
	return js.Value(el).Call("togglePopover", force)
}


func (el Element) AddEventListener2K-(this: HTMLElement, ev: HTMLElementEventMap[K]) => any(typeName any, listener any) js.Value {
	return js.Value(el).Call("addEventListener", typeName, listener)
}


func (el Element) AddEventListener2StringEventListenerOrEventListenerObject(typeName string, listener any) js.Value {
	return js.Value(el).Call("addEventListener", typeName, listener)
}


func (el Element) AddEventListener3K-(this: HTMLElement, ev: HTMLElementEventMap[K]) => anyBoolean(typeName any, listener any, options bool) js.Value {
	return js.Value(el).Call("addEventListener", typeName, listener, options)
}


func (el Element) AddEventListener3StringEventListenerOrEventListenerObjectBoolean(typeName string, listener any, options bool) js.Value {
	return js.Value(el).Call("addEventListener", typeName, listener, options)
}


func (el Element) RemoveEventListener2K-(this: HTMLElement, ev: HTMLElementEventMap[K]) => any(typeName any, listener any) js.Value {
	return js.Value(el).Call("removeEventListener", typeName, listener)
}


func (el Element) RemoveEventListener2StringEventListenerOrEventListenerObject(typeName string, listener any) js.Value {
	return js.Value(el).Call("removeEventListener", typeName, listener)
}


func (el Element) RemoveEventListener3K-(this: HTMLElement, ev: HTMLElementEventMap[K]) => anyBoolean(typeName any, listener any, options bool) js.Value {
	return js.Value(el).Call("removeEventListener", typeName, listener, options)
}


func (el Element) RemoveEventListener3StringEventListenerOrEventListenerObjectBoolean(typeName string, listener any, options bool) js.Value {
	return js.Value(el).Call("removeEventListener", typeName, listener, options)
}


func (el Element) AttachShadow(init any) js.Value {
	return js.Value(el).Call("attachShadow", init)
}


func (el Element) CheckVisibility0() js.Value {
	return js.Value(el).Call("checkVisibility")
}


func (el Element) CheckVisibility1(options any) js.Value {
	return js.Value(el).Call("checkVisibility", options)
}


func (el Element) Closest1K(selector any) js.Value {
	return js.Value(el).Call("closest", selector)
}


func (el Element) Closest1String(selectors string) js.Value {
	return js.Value(el).Call("closest", selectors)
}


func (el Element) ComputedStyleMap() js.Value {
	return js.Value(el).Call("computedStyleMap")
}


func (el Element) GetAttribute(qualifiedName string) js.Value {
	return js.Value(el).Call("getAttribute", qualifiedName)
}


func (el Element) GetAttributeNS(namespace, localName string) js.Value {
	return js.Value(el).Call("getAttributeNS", namespace, localName)
}


func (el Element) GetAttributeNames() js.Value {
	return js.Value(el).Call("getAttributeNames")
}


func (el Element) GetAttributeNode(qualifiedName string) js.Value {
	return js.Value(el).Call("getAttributeNode", qualifiedName)
}


func (el Element) GetAttributeNodeNS(namespace, localName string) js.Value {
	return js.Value(el).Call("getAttributeNodeNS", namespace, localName)
}


func (el Element) GetBoundingClientRect() js.Value {
	return js.Value(el).Call("getBoundingClientRect")
}


func (el Element) GetClientRects() js.Value {
	return js.Value(el).Call("getClientRects")
}


func (el Element) GetElementsByClassName(classNames string) js.Value {
	return js.Value(el).Call("getElementsByClassName", classNames)
}


func (el Element) GetElementsByTagName1K(qualifiedName any) js.Value {
	return js.Value(el).Call("getElementsByTagName", qualifiedName)
}


func (el Element) GetElementsByTagName1String(qualifiedName string) js.Value {
	return js.Value(el).Call("getElementsByTagName", qualifiedName)
}


func (el Element) GetElementsByTagNameNS(namespace, localName string) js.Value {
	return js.Value(el).Call("getElementsByTagNameNS", namespace, localName)
}


func (el Element) HasAttribute(qualifiedName string) js.Value {
	return js.Value(el).Call("hasAttribute", qualifiedName)
}


func (el Element) HasAttributeNS(namespace, localName string) js.Value {
	return js.Value(el).Call("hasAttributeNS", namespace, localName)
}


func (el Element) HasAttributes() js.Value {
	return js.Value(el).Call("hasAttributes")
}


func (el Element) HasPointerCapture(pointerId float64) js.Value {
	return js.Value(el).Call("hasPointerCapture", pointerId)
}


func (el Element) InsertAdjacentElement(where InsertPosition, element any) js.Value {
	return js.Value(el).Call("insertAdjacentElement", where, element)
}


func (el Element) InsertAdjacentHTML(position InsertPosition, text string) js.Value {
	return js.Value(el).Call("insertAdjacentHTML", position, text)
}


func (el Element) InsertAdjacentText(where InsertPosition, data string) js.Value {
	return js.Value(el).Call("insertAdjacentText", where, data)
}


func (el Element) Matches(selectors string) js.Value {
	return js.Value(el).Call("matches", selectors)
}


func (el Element) ReleasePointerCapture(pointerId float64) js.Value {
	return js.Value(el).Call("releasePointerCapture", pointerId)
}


func (el Element) RemoveAttribute(qualifiedName string) js.Value {
	return js.Value(el).Call("removeAttribute", qualifiedName)
}


func (el Element) RemoveAttributeNS(namespace, localName string) js.Value {
	return js.Value(el).Call("removeAttributeNS", namespace, localName)
}


func (el Element) RemoveAttributeNode(attr any) js.Value {
	return js.Value(el).Call("removeAttributeNode", attr)
}


func (el Element) RequestFullscreen0() js.Value {
	return js.Value(el).Call("requestFullscreen")
}


func (el Element) RequestFullscreen1(options any) js.Value {
	return js.Value(el).Call("requestFullscreen", options)
}


func (el Element) RequestPointerLock() js.Value {
	return js.Value(el).Call("requestPointerLock")
}


func (el Element) Scroll0() js.Value {
	return js.Value(el).Call("scroll")
}


func (el Element) Scroll1(options any) js.Value {
	return js.Value(el).Call("scroll", options)
}


func (el Element) Scroll2(x, y float64) js.Value {
	return js.Value(el).Call("scroll", x, y)
}


func (el Element) ScrollBy0() js.Value {
	return js.Value(el).Call("scrollBy")
}


func (el Element) ScrollBy1(options any) js.Value {
	return js.Value(el).Call("scrollBy", options)
}


func (el Element) ScrollBy2(x, y float64) js.Value {
	return js.Value(el).Call("scrollBy", x, y)
}


func (el Element) ScrollIntoView0() js.Value {
	return js.Value(el).Call("scrollIntoView")
}


func (el Element) ScrollIntoView1(arg bool) js.Value {
	return js.Value(el).Call("scrollIntoView", arg)
}


func (el Element) ScrollTo0() js.Value {
	return js.Value(el).Call("scrollTo")
}


func (el Element) ScrollTo1(options any) js.Value {
	return js.Value(el).Call("scrollTo", options)
}


func (el Element) ScrollTo2(x, y float64) js.Value {
	return js.Value(el).Call("scrollTo", x, y)
}


func (el Element) SetAttribute(qualifiedName, value string) js.Value {
	return js.Value(el).Call("setAttribute", qualifiedName, value)
}


func (el Element) SetAttributeNS(namespace, qualifiedName, value string) js.Value {
	return js.Value(el).Call("setAttributeNS", namespace, qualifiedName, value)
}


func (el Element) SetAttributeNode(attr any) js.Value {
	return js.Value(el).Call("setAttributeNode", attr)
}


func (el Element) SetAttributeNodeNS(attr any) js.Value {
	return js.Value(el).Call("setAttributeNodeNS", attr)
}


func (el Element) SetPointerCapture(pointerId float64) js.Value {
	return js.Value(el).Call("setPointerCapture", pointerId)
}


func (el Element) ToggleAttribute1(qualifiedName string) js.Value {
	return js.Value(el).Call("toggleAttribute", qualifiedName)
}


func (el Element) ToggleAttribute2(qualifiedName string, force bool) js.Value {
	return js.Value(el).Call("toggleAttribute", qualifiedName, force)
}


func (el Element) WebkitMatchesSelector(selectors string) js.Value {
	return js.Value(el).Call("webkitMatchesSelector", selectors)
}


func (el Element) AppendChild(node any) js.Value {
	return js.Value(el).Call("appendChild", node)
}


func (el Element) CloneNode0() js.Value {
	return js.Value(el).Call("cloneNode")
}


func (el Element) CloneNode1(deep bool) js.Value {
	return js.Value(el).Call("cloneNode", deep)
}


func (el Element) CompareDocumentPosition(other any) js.Value {
	return js.Value(el).Call("compareDocumentPosition", other)
}


func (el Element) Contains(other any) js.Value {
	return js.Value(el).Call("contains", other)
}


func (el Element) GetRootNode0() js.Value {
	return js.Value(el).Call("getRootNode")
}


func (el Element) GetRootNode1(options any) js.Value {
	return js.Value(el).Call("getRootNode", options)
}


func (el Element) HasChildNodes() js.Value {
	return js.Value(el).Call("hasChildNodes")
}


func (el Element) InsertBefore(node any, child any) js.Value {
	return js.Value(el).Call("insertBefore", node, child)
}


func (el Element) IsDefaultNamespace(namespace string) js.Value {
	return js.Value(el).Call("isDefaultNamespace", namespace)
}


func (el Element) IsEqualNode(otherNode any) js.Value {
	return js.Value(el).Call("isEqualNode", otherNode)
}


func (el Element) IsSameNode(otherNode any) js.Value {
	return js.Value(el).Call("isSameNode", otherNode)
}


func (el Element) LookupNamespaceURI(prefix string) js.Value {
	return js.Value(el).Call("lookupNamespaceURI", prefix)
}


func (el Element) LookupPrefix(namespace string) js.Value {
	return js.Value(el).Call("lookupPrefix", namespace)
}


func (el Element) Normalize() js.Value {
	return js.Value(el).Call("normalize")
}


func (el Element) RemoveChild(child any) js.Value {
	return js.Value(el).Call("removeChild", child)
}


func (el Element) ReplaceChild(node any, child any) js.Value {
	return js.Value(el).Call("replaceChild", node, child)
}


func (el Element) DispatchEvent(event any) js.Value {
	return js.Value(el).Call("dispatchEvent", event)
}


func (el Element) Animate1(keyframes any) js.Value {
	return js.Value(el).Call("animate", keyframes)
}


func (el Element) Animate2(keyframes any, options float64) js.Value {
	return js.Value(el).Call("animate", keyframes, options)
}


func (el Element) GetAnimations0() js.Value {
	return js.Value(el).Call("getAnimations")
}


func (el Element) GetAnimations1(options any) js.Value {
	return js.Value(el).Call("getAnimations", options)
}


func (el Element) After(nodes string) js.Value {
	return js.Value(el).Call("after", nodes)
}


func (el Element) Before(nodes string) js.Value {
	return js.Value(el).Call("before", nodes)
}


func (el Element) Remove() js.Value {
	return js.Value(el).Call("remove")
}


func (el Element) ReplaceWith(nodes string) js.Value {
	return js.Value(el).Call("replaceWith", nodes)
}


func (el Element) Append(nodes string) js.Value {
	return js.Value(el).Call("append", nodes)
}


func (el Element) Prepend(nodes string) js.Value {
	return js.Value(el).Call("prepend", nodes)
}


func (el Element) QuerySelector1K(selectors any) js.Value {
	return js.Value(el).Call("querySelector", selectors)
}


func (el Element) QuerySelector1String(selectors string) js.Value {
	return js.Value(el).Call("querySelector", selectors)
}


func (el Element) QuerySelectorAll1K(selectors any) js.Value {
	return js.Value(el).Call("querySelectorAll", selectors)
}


func (el Element) QuerySelectorAll1String(selectors string) js.Value {
	return js.Value(el).Call("querySelectorAll", selectors)
}


func (el Element) ReplaceChildren(nodes string) js.Value {
	return js.Value(el).Call("replaceChildren", nodes)
}


func (el Element) Blur() js.Value {
	return js.Value(el).Call("blur")
}


func (el Element) Focus0() js.Value {
	return js.Value(el).Call("focus")
}


func (el Element) Focus1(options any) js.Value {
	return js.Value(el).Call("focus", options)
}



type InsertPosition = string
const (
  InsertPositionAfterbegin InsertPosition = "afterbegin"
  InsertPositionAfterend InsertPosition = "afterend"
  InsertPositionBeforebegin InsertPosition = "beforebegin"
  InsertPositionBeforeend InsertPosition = "beforeend"
)



