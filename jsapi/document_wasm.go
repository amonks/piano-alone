package jsapi

import "syscall/js"

type Document js.Value


func (doc Document) GetURL() string {
	return js.Value(doc).Get("URL").String()
}

func (doc Document) SetURL(val string) {
	js.Value(doc).Set("URL", val)
}


func (doc Document) GetAlinkColor() string {
	return js.Value(doc).Get("alinkColor").String()
}

func (doc Document) SetAlinkColor(val string) {
	js.Value(doc).Set("alinkColor", val)
}


func (doc Document) GetAll() js.Value {
	return js.Value(doc).Get("all")
}

func (doc Document) SetAll(val any) {
	js.Value(doc).Set("all", val)
}


func (doc Document) GetAnchors() js.Value {
	return js.Value(doc).Get("anchors")
}

func (doc Document) SetAnchors(val any) {
	js.Value(doc).Set("anchors", val)
}


func (doc Document) GetApplets() js.Value {
	return js.Value(doc).Get("applets")
}

func (doc Document) SetApplets(val any) {
	js.Value(doc).Set("applets", val)
}


func (doc Document) GetBgColor() string {
	return js.Value(doc).Get("bgColor").String()
}

func (doc Document) SetBgColor(val string) {
	js.Value(doc).Set("bgColor", val)
}


func (doc Document) GetBody() js.Value {
	return js.Value(doc).Get("body")
}

func (doc Document) SetBody(val any) {
	js.Value(doc).Set("body", val)
}


func (doc Document) GetCharacterSet() string {
	return js.Value(doc).Get("characterSet").String()
}

func (doc Document) SetCharacterSet(val string) {
	js.Value(doc).Set("characterSet", val)
}


func (doc Document) GetCharset() string {
	return js.Value(doc).Get("charset").String()
}

func (doc Document) SetCharset(val string) {
	js.Value(doc).Set("charset", val)
}


func (doc Document) GetCompatMode() string {
	return js.Value(doc).Get("compatMode").String()
}

func (doc Document) SetCompatMode(val string) {
	js.Value(doc).Set("compatMode", val)
}


func (doc Document) GetContentType() string {
	return js.Value(doc).Get("contentType").String()
}

func (doc Document) SetContentType(val string) {
	js.Value(doc).Set("contentType", val)
}


func (doc Document) GetCookie() string {
	return js.Value(doc).Get("cookie").String()
}

func (doc Document) SetCookie(val string) {
	js.Value(doc).Set("cookie", val)
}


func (doc Document) GetCurrentScript() js.Value {
	return js.Value(doc).Get("currentScript")
}

func (doc Document) SetCurrentScript(val any) {
	js.Value(doc).Set("currentScript", val)
}


func (doc Document) GetDefaultView() js.Value {
	return js.Value(doc).Get("defaultView")
}

func (doc Document) SetDefaultView(val any) {
	js.Value(doc).Set("defaultView", val)
}


func (doc Document) GetDesignMode() string {
	return js.Value(doc).Get("designMode").String()
}

func (doc Document) SetDesignMode(val string) {
	js.Value(doc).Set("designMode", val)
}


func (doc Document) GetDir() string {
	return js.Value(doc).Get("dir").String()
}

func (doc Document) SetDir(val string) {
	js.Value(doc).Set("dir", val)
}


func (doc Document) GetDoctype() js.Value {
	return js.Value(doc).Get("doctype")
}

func (doc Document) SetDoctype(val any) {
	js.Value(doc).Set("doctype", val)
}


func (doc Document) GetDocumentElement() js.Value {
	return js.Value(doc).Get("documentElement")
}

func (doc Document) SetDocumentElement(val any) {
	js.Value(doc).Set("documentElement", val)
}


func (doc Document) GetDocumentURI() string {
	return js.Value(doc).Get("documentURI").String()
}

func (doc Document) SetDocumentURI(val string) {
	js.Value(doc).Set("documentURI", val)
}


func (doc Document) GetDomain() string {
	return js.Value(doc).Get("domain").String()
}

func (doc Document) SetDomain(val string) {
	js.Value(doc).Set("domain", val)
}


func (doc Document) GetEmbeds() js.Value {
	return js.Value(doc).Get("embeds")
}

func (doc Document) SetEmbeds(val any) {
	js.Value(doc).Set("embeds", val)
}


func (doc Document) GetFgColor() string {
	return js.Value(doc).Get("fgColor").String()
}

func (doc Document) SetFgColor(val string) {
	js.Value(doc).Set("fgColor", val)
}


func (doc Document) GetForms() js.Value {
	return js.Value(doc).Get("forms")
}

func (doc Document) SetForms(val any) {
	js.Value(doc).Set("forms", val)
}


func (doc Document) GetFullscreen() bool {
	return js.Value(doc).Get("fullscreen").Bool()
}

func (doc Document) SetFullscreen(val bool) {
	js.Value(doc).Set("fullscreen", val)
}


func (doc Document) GetFullscreenEnabled() bool {
	return js.Value(doc).Get("fullscreenEnabled").Bool()
}

func (doc Document) SetFullscreenEnabled(val bool) {
	js.Value(doc).Set("fullscreenEnabled", val)
}


func (doc Document) GetHead() js.Value {
	return js.Value(doc).Get("head")
}

func (doc Document) SetHead(val any) {
	js.Value(doc).Set("head", val)
}


func (doc Document) GetHidden() bool {
	return js.Value(doc).Get("hidden").Bool()
}

func (doc Document) SetHidden(val bool) {
	js.Value(doc).Set("hidden", val)
}


func (doc Document) GetImages() js.Value {
	return js.Value(doc).Get("images")
}

func (doc Document) SetImages(val any) {
	js.Value(doc).Set("images", val)
}


func (doc Document) GetImplementation() js.Value {
	return js.Value(doc).Get("implementation")
}

func (doc Document) SetImplementation(val any) {
	js.Value(doc).Set("implementation", val)
}


func (doc Document) GetInputEncoding() string {
	return js.Value(doc).Get("inputEncoding").String()
}

func (doc Document) SetInputEncoding(val string) {
	js.Value(doc).Set("inputEncoding", val)
}


func (doc Document) GetLastModified() string {
	return js.Value(doc).Get("lastModified").String()
}

func (doc Document) SetLastModified(val string) {
	js.Value(doc).Set("lastModified", val)
}


func (doc Document) GetLinkColor() string {
	return js.Value(doc).Get("linkColor").String()
}

func (doc Document) SetLinkColor(val string) {
	js.Value(doc).Set("linkColor", val)
}


func (doc Document) GetLinks() js.Value {
	return js.Value(doc).Get("links")
}

func (doc Document) SetLinks(val any) {
	js.Value(doc).Set("links", val)
}


func (doc Document) GetOnfullscreenchange() js.Value {
	return js.Value(doc).Get("onfullscreenchange")
}

func (doc Document) SetOnfullscreenchange(val any) {
	js.Value(doc).Set("onfullscreenchange", val)
}


func (doc Document) GetOnfullscreenerror() js.Value {
	return js.Value(doc).Get("onfullscreenerror")
}

func (doc Document) SetOnfullscreenerror(val any) {
	js.Value(doc).Set("onfullscreenerror", val)
}


func (doc Document) GetOnpointerlockchange() js.Value {
	return js.Value(doc).Get("onpointerlockchange")
}

func (doc Document) SetOnpointerlockchange(val any) {
	js.Value(doc).Set("onpointerlockchange", val)
}


func (doc Document) GetOnpointerlockerror() js.Value {
	return js.Value(doc).Get("onpointerlockerror")
}

func (doc Document) SetOnpointerlockerror(val any) {
	js.Value(doc).Set("onpointerlockerror", val)
}


func (doc Document) GetOnreadystatechange() js.Value {
	return js.Value(doc).Get("onreadystatechange")
}

func (doc Document) SetOnreadystatechange(val any) {
	js.Value(doc).Set("onreadystatechange", val)
}


func (doc Document) GetOnvisibilitychange() js.Value {
	return js.Value(doc).Get("onvisibilitychange")
}

func (doc Document) SetOnvisibilitychange(val any) {
	js.Value(doc).Set("onvisibilitychange", val)
}


func (doc Document) GetOwnerDocument() js.Value {
	return js.Value(doc).Get("ownerDocument")
}

func (doc Document) SetOwnerDocument(val any) {
	js.Value(doc).Set("ownerDocument", val)
}


func (doc Document) GetPictureInPictureEnabled() bool {
	return js.Value(doc).Get("pictureInPictureEnabled").Bool()
}

func (doc Document) SetPictureInPictureEnabled(val bool) {
	js.Value(doc).Set("pictureInPictureEnabled", val)
}


func (doc Document) GetPlugins() js.Value {
	return js.Value(doc).Get("plugins")
}

func (doc Document) SetPlugins(val any) {
	js.Value(doc).Set("plugins", val)
}


func (doc Document) GetReadyState() DocumentReadyState {
	return js.Value(doc).Get("readyState").String()
}

func (doc Document) SetReadyState(val DocumentReadyState) {
	js.Value(doc).Set("readyState", val)
}


func (doc Document) GetReferrer() string {
	return js.Value(doc).Get("referrer").String()
}

func (doc Document) SetReferrer(val string) {
	js.Value(doc).Set("referrer", val)
}


func (doc Document) GetRootElement() js.Value {
	return js.Value(doc).Get("rootElement")
}

func (doc Document) SetRootElement(val any) {
	js.Value(doc).Set("rootElement", val)
}


func (doc Document) GetScripts() js.Value {
	return js.Value(doc).Get("scripts")
}

func (doc Document) SetScripts(val any) {
	js.Value(doc).Set("scripts", val)
}


func (doc Document) GetScrollingElement() js.Value {
	return js.Value(doc).Get("scrollingElement")
}

func (doc Document) SetScrollingElement(val any) {
	js.Value(doc).Set("scrollingElement", val)
}


func (doc Document) GetTimeline() js.Value {
	return js.Value(doc).Get("timeline")
}

func (doc Document) SetTimeline(val any) {
	js.Value(doc).Set("timeline", val)
}


func (doc Document) GetTitle() string {
	return js.Value(doc).Get("title").String()
}

func (doc Document) SetTitle(val string) {
	js.Value(doc).Set("title", val)
}


func (doc Document) GetVisibilityState() DocumentVisibilityState {
	return js.Value(doc).Get("visibilityState").String()
}

func (doc Document) SetVisibilityState(val DocumentVisibilityState) {
	js.Value(doc).Set("visibilityState", val)
}


func (doc Document) GetVlinkColor() string {
	return js.Value(doc).Get("vlinkColor").String()
}

func (doc Document) SetVlinkColor(val string) {
	js.Value(doc).Set("vlinkColor", val)
}


func (doc Document) GetBaseURI() string {
	return js.Value(doc).Get("baseURI").String()
}

func (doc Document) SetBaseURI(val string) {
	js.Value(doc).Set("baseURI", val)
}


func (doc Document) GetChildNodes() js.Value {
	return js.Value(doc).Get("childNodes")
}

func (doc Document) SetChildNodes(val any) {
	js.Value(doc).Set("childNodes", val)
}


func (doc Document) GetFirstChild() js.Value {
	return js.Value(doc).Get("firstChild")
}

func (doc Document) SetFirstChild(val any) {
	js.Value(doc).Set("firstChild", val)
}


func (doc Document) GetIsConnected() bool {
	return js.Value(doc).Get("isConnected").Bool()
}

func (doc Document) SetIsConnected(val bool) {
	js.Value(doc).Set("isConnected", val)
}


func (doc Document) GetLastChild() js.Value {
	return js.Value(doc).Get("lastChild")
}

func (doc Document) SetLastChild(val any) {
	js.Value(doc).Set("lastChild", val)
}


func (doc Document) GetNextSibling() js.Value {
	return js.Value(doc).Get("nextSibling")
}

func (doc Document) SetNextSibling(val any) {
	js.Value(doc).Set("nextSibling", val)
}


func (doc Document) GetNodeName() string {
	return js.Value(doc).Get("nodeName").String()
}

func (doc Document) SetNodeName(val string) {
	js.Value(doc).Set("nodeName", val)
}


func (doc Document) GetNodeType() float64 {
	return js.Value(doc).Get("nodeType").Float()
}

func (doc Document) SetNodeType(val float64) {
	js.Value(doc).Set("nodeType", val)
}


func (doc Document) GetNodeValue() string {
	return js.Value(doc).Get("nodeValue").String()
}

func (doc Document) SetNodeValue(val string) {
	js.Value(doc).Set("nodeValue", val)
}


func (doc Document) GetParentElement() js.Value {
	return js.Value(doc).Get("parentElement")
}

func (doc Document) SetParentElement(val any) {
	js.Value(doc).Set("parentElement", val)
}


func (doc Document) GetParentNode() js.Value {
	return js.Value(doc).Get("parentNode")
}

func (doc Document) SetParentNode(val any) {
	js.Value(doc).Set("parentNode", val)
}


func (doc Document) GetPreviousSibling() js.Value {
	return js.Value(doc).Get("previousSibling")
}

func (doc Document) SetPreviousSibling(val any) {
	js.Value(doc).Set("previousSibling", val)
}


func (doc Document) GetTextContent() string {
	return js.Value(doc).Get("textContent").String()
}

func (doc Document) SetTextContent(val string) {
	js.Value(doc).Set("textContent", val)
}


func (doc Document) GetELEMENT_NODE() js.Value {
	return js.Value(doc).Get("ELEMENT_NODE")
}

func (doc Document) SetELEMENT_NODE(val any) {
	js.Value(doc).Set("ELEMENT_NODE", val)
}


func (doc Document) GetATTRIBUTE_NODE() js.Value {
	return js.Value(doc).Get("ATTRIBUTE_NODE")
}

func (doc Document) SetATTRIBUTE_NODE(val any) {
	js.Value(doc).Set("ATTRIBUTE_NODE", val)
}


func (doc Document) GetTEXT_NODE() js.Value {
	return js.Value(doc).Get("TEXT_NODE")
}

func (doc Document) SetTEXT_NODE(val any) {
	js.Value(doc).Set("TEXT_NODE", val)
}


func (doc Document) GetCDATA_SECTION_NODE() js.Value {
	return js.Value(doc).Get("CDATA_SECTION_NODE")
}

func (doc Document) SetCDATA_SECTION_NODE(val any) {
	js.Value(doc).Set("CDATA_SECTION_NODE", val)
}


func (doc Document) GetENTITY_REFERENCE_NODE() js.Value {
	return js.Value(doc).Get("ENTITY_REFERENCE_NODE")
}

func (doc Document) SetENTITY_REFERENCE_NODE(val any) {
	js.Value(doc).Set("ENTITY_REFERENCE_NODE", val)
}


func (doc Document) GetENTITY_NODE() js.Value {
	return js.Value(doc).Get("ENTITY_NODE")
}

func (doc Document) SetENTITY_NODE(val any) {
	js.Value(doc).Set("ENTITY_NODE", val)
}


func (doc Document) GetPROCESSING_INSTRUCTION_NODE() js.Value {
	return js.Value(doc).Get("PROCESSING_INSTRUCTION_NODE")
}

func (doc Document) SetPROCESSING_INSTRUCTION_NODE(val any) {
	js.Value(doc).Set("PROCESSING_INSTRUCTION_NODE", val)
}


func (doc Document) GetCOMMENT_NODE() js.Value {
	return js.Value(doc).Get("COMMENT_NODE")
}

func (doc Document) SetCOMMENT_NODE(val any) {
	js.Value(doc).Set("COMMENT_NODE", val)
}


func (doc Document) GetDOCUMENT_NODE() js.Value {
	return js.Value(doc).Get("DOCUMENT_NODE")
}

func (doc Document) SetDOCUMENT_NODE(val any) {
	js.Value(doc).Set("DOCUMENT_NODE", val)
}


func (doc Document) GetDOCUMENT_TYPE_NODE() js.Value {
	return js.Value(doc).Get("DOCUMENT_TYPE_NODE")
}

func (doc Document) SetDOCUMENT_TYPE_NODE(val any) {
	js.Value(doc).Set("DOCUMENT_TYPE_NODE", val)
}


func (doc Document) GetDOCUMENT_FRAGMENT_NODE() js.Value {
	return js.Value(doc).Get("DOCUMENT_FRAGMENT_NODE")
}

func (doc Document) SetDOCUMENT_FRAGMENT_NODE(val any) {
	js.Value(doc).Set("DOCUMENT_FRAGMENT_NODE", val)
}


func (doc Document) GetNOTATION_NODE() js.Value {
	return js.Value(doc).Get("NOTATION_NODE")
}

func (doc Document) SetNOTATION_NODE(val any) {
	js.Value(doc).Set("NOTATION_NODE", val)
}


func (doc Document) GetDOCUMENT_POSITION_DISCONNECTED() js.Value {
	return js.Value(doc).Get("DOCUMENT_POSITION_DISCONNECTED")
}

func (doc Document) SetDOCUMENT_POSITION_DISCONNECTED(val any) {
	js.Value(doc).Set("DOCUMENT_POSITION_DISCONNECTED", val)
}


func (doc Document) GetDOCUMENT_POSITION_PRECEDING() js.Value {
	return js.Value(doc).Get("DOCUMENT_POSITION_PRECEDING")
}

func (doc Document) SetDOCUMENT_POSITION_PRECEDING(val any) {
	js.Value(doc).Set("DOCUMENT_POSITION_PRECEDING", val)
}


func (doc Document) GetDOCUMENT_POSITION_FOLLOWING() js.Value {
	return js.Value(doc).Get("DOCUMENT_POSITION_FOLLOWING")
}

func (doc Document) SetDOCUMENT_POSITION_FOLLOWING(val any) {
	js.Value(doc).Set("DOCUMENT_POSITION_FOLLOWING", val)
}


func (doc Document) GetDOCUMENT_POSITION_CONTAINS() js.Value {
	return js.Value(doc).Get("DOCUMENT_POSITION_CONTAINS")
}

func (doc Document) SetDOCUMENT_POSITION_CONTAINS(val any) {
	js.Value(doc).Set("DOCUMENT_POSITION_CONTAINS", val)
}


func (doc Document) GetDOCUMENT_POSITION_CONTAINED_BY() js.Value {
	return js.Value(doc).Get("DOCUMENT_POSITION_CONTAINED_BY")
}

func (doc Document) SetDOCUMENT_POSITION_CONTAINED_BY(val any) {
	js.Value(doc).Set("DOCUMENT_POSITION_CONTAINED_BY", val)
}


func (doc Document) GetDOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC() js.Value {
	return js.Value(doc).Get("DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC")
}

func (doc Document) SetDOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC(val any) {
	js.Value(doc).Set("DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC", val)
}


func (doc Document) GetActiveElement() js.Value {
	return js.Value(doc).Get("activeElement")
}

func (doc Document) SetActiveElement(val any) {
	js.Value(doc).Set("activeElement", val)
}


func (doc Document) GetAdoptedStyleSheets() js.Value {
	return js.Value(doc).Get("adoptedStyleSheets")
}

func (doc Document) SetAdoptedStyleSheets(val any) {
	js.Value(doc).Set("adoptedStyleSheets", val)
}


func (doc Document) GetFullscreenElement() js.Value {
	return js.Value(doc).Get("fullscreenElement")
}

func (doc Document) SetFullscreenElement(val any) {
	js.Value(doc).Set("fullscreenElement", val)
}


func (doc Document) GetPictureInPictureElement() js.Value {
	return js.Value(doc).Get("pictureInPictureElement")
}

func (doc Document) SetPictureInPictureElement(val any) {
	js.Value(doc).Set("pictureInPictureElement", val)
}


func (doc Document) GetPointerLockElement() js.Value {
	return js.Value(doc).Get("pointerLockElement")
}

func (doc Document) SetPointerLockElement(val any) {
	js.Value(doc).Set("pointerLockElement", val)
}


func (doc Document) GetStyleSheets() js.Value {
	return js.Value(doc).Get("styleSheets")
}

func (doc Document) SetStyleSheets(val any) {
	js.Value(doc).Set("styleSheets", val)
}


func (doc Document) GetFonts() js.Value {
	return js.Value(doc).Get("fonts")
}

func (doc Document) SetFonts(val any) {
	js.Value(doc).Set("fonts", val)
}


func (doc Document) GetOnabort() js.Value {
	return js.Value(doc).Get("onabort")
}

func (doc Document) SetOnabort(val any) {
	js.Value(doc).Set("onabort", val)
}


func (doc Document) GetOnanimationcancel() js.Value {
	return js.Value(doc).Get("onanimationcancel")
}

func (doc Document) SetOnanimationcancel(val any) {
	js.Value(doc).Set("onanimationcancel", val)
}


func (doc Document) GetOnanimationend() js.Value {
	return js.Value(doc).Get("onanimationend")
}

func (doc Document) SetOnanimationend(val any) {
	js.Value(doc).Set("onanimationend", val)
}


func (doc Document) GetOnanimationiteration() js.Value {
	return js.Value(doc).Get("onanimationiteration")
}

func (doc Document) SetOnanimationiteration(val any) {
	js.Value(doc).Set("onanimationiteration", val)
}


func (doc Document) GetOnanimationstart() js.Value {
	return js.Value(doc).Get("onanimationstart")
}

func (doc Document) SetOnanimationstart(val any) {
	js.Value(doc).Set("onanimationstart", val)
}


func (doc Document) GetOnauxclick() js.Value {
	return js.Value(doc).Get("onauxclick")
}

func (doc Document) SetOnauxclick(val any) {
	js.Value(doc).Set("onauxclick", val)
}


func (doc Document) GetOnbeforeinput() js.Value {
	return js.Value(doc).Get("onbeforeinput")
}

func (doc Document) SetOnbeforeinput(val any) {
	js.Value(doc).Set("onbeforeinput", val)
}


func (doc Document) GetOnblur() js.Value {
	return js.Value(doc).Get("onblur")
}

func (doc Document) SetOnblur(val any) {
	js.Value(doc).Set("onblur", val)
}


func (doc Document) GetOncancel() js.Value {
	return js.Value(doc).Get("oncancel")
}

func (doc Document) SetOncancel(val any) {
	js.Value(doc).Set("oncancel", val)
}


func (doc Document) GetOncanplay() js.Value {
	return js.Value(doc).Get("oncanplay")
}

func (doc Document) SetOncanplay(val any) {
	js.Value(doc).Set("oncanplay", val)
}


func (doc Document) GetOncanplaythrough() js.Value {
	return js.Value(doc).Get("oncanplaythrough")
}

func (doc Document) SetOncanplaythrough(val any) {
	js.Value(doc).Set("oncanplaythrough", val)
}


func (doc Document) GetOnchange() js.Value {
	return js.Value(doc).Get("onchange")
}

func (doc Document) SetOnchange(val any) {
	js.Value(doc).Set("onchange", val)
}


func (doc Document) GetOnclick() js.Value {
	return js.Value(doc).Get("onclick")
}

func (doc Document) SetOnclick(val any) {
	js.Value(doc).Set("onclick", val)
}


func (doc Document) GetOnclose() js.Value {
	return js.Value(doc).Get("onclose")
}

func (doc Document) SetOnclose(val any) {
	js.Value(doc).Set("onclose", val)
}


func (doc Document) GetOncontextmenu() js.Value {
	return js.Value(doc).Get("oncontextmenu")
}

func (doc Document) SetOncontextmenu(val any) {
	js.Value(doc).Set("oncontextmenu", val)
}


func (doc Document) GetOncopy() js.Value {
	return js.Value(doc).Get("oncopy")
}

func (doc Document) SetOncopy(val any) {
	js.Value(doc).Set("oncopy", val)
}


func (doc Document) GetOncuechange() js.Value {
	return js.Value(doc).Get("oncuechange")
}

func (doc Document) SetOncuechange(val any) {
	js.Value(doc).Set("oncuechange", val)
}


func (doc Document) GetOncut() js.Value {
	return js.Value(doc).Get("oncut")
}

func (doc Document) SetOncut(val any) {
	js.Value(doc).Set("oncut", val)
}


func (doc Document) GetOndblclick() js.Value {
	return js.Value(doc).Get("ondblclick")
}

func (doc Document) SetOndblclick(val any) {
	js.Value(doc).Set("ondblclick", val)
}


func (doc Document) GetOndrag() js.Value {
	return js.Value(doc).Get("ondrag")
}

func (doc Document) SetOndrag(val any) {
	js.Value(doc).Set("ondrag", val)
}


func (doc Document) GetOndragend() js.Value {
	return js.Value(doc).Get("ondragend")
}

func (doc Document) SetOndragend(val any) {
	js.Value(doc).Set("ondragend", val)
}


func (doc Document) GetOndragenter() js.Value {
	return js.Value(doc).Get("ondragenter")
}

func (doc Document) SetOndragenter(val any) {
	js.Value(doc).Set("ondragenter", val)
}


func (doc Document) GetOndragleave() js.Value {
	return js.Value(doc).Get("ondragleave")
}

func (doc Document) SetOndragleave(val any) {
	js.Value(doc).Set("ondragleave", val)
}


func (doc Document) GetOndragover() js.Value {
	return js.Value(doc).Get("ondragover")
}

func (doc Document) SetOndragover(val any) {
	js.Value(doc).Set("ondragover", val)
}


func (doc Document) GetOndragstart() js.Value {
	return js.Value(doc).Get("ondragstart")
}

func (doc Document) SetOndragstart(val any) {
	js.Value(doc).Set("ondragstart", val)
}


func (doc Document) GetOndrop() js.Value {
	return js.Value(doc).Get("ondrop")
}

func (doc Document) SetOndrop(val any) {
	js.Value(doc).Set("ondrop", val)
}


func (doc Document) GetOndurationchange() js.Value {
	return js.Value(doc).Get("ondurationchange")
}

func (doc Document) SetOndurationchange(val any) {
	js.Value(doc).Set("ondurationchange", val)
}


func (doc Document) GetOnemptied() js.Value {
	return js.Value(doc).Get("onemptied")
}

func (doc Document) SetOnemptied(val any) {
	js.Value(doc).Set("onemptied", val)
}


func (doc Document) GetOnended() js.Value {
	return js.Value(doc).Get("onended")
}

func (doc Document) SetOnended(val any) {
	js.Value(doc).Set("onended", val)
}


func (doc Document) GetOnerror() js.Value {
	return js.Value(doc).Get("onerror")
}

func (doc Document) SetOnerror(val any) {
	js.Value(doc).Set("onerror", val)
}


func (doc Document) GetOnfocus() js.Value {
	return js.Value(doc).Get("onfocus")
}

func (doc Document) SetOnfocus(val any) {
	js.Value(doc).Set("onfocus", val)
}


func (doc Document) GetOnformdata() js.Value {
	return js.Value(doc).Get("onformdata")
}

func (doc Document) SetOnformdata(val any) {
	js.Value(doc).Set("onformdata", val)
}


func (doc Document) GetOngotpointercapture() js.Value {
	return js.Value(doc).Get("ongotpointercapture")
}

func (doc Document) SetOngotpointercapture(val any) {
	js.Value(doc).Set("ongotpointercapture", val)
}


func (doc Document) GetOninput() js.Value {
	return js.Value(doc).Get("oninput")
}

func (doc Document) SetOninput(val any) {
	js.Value(doc).Set("oninput", val)
}


func (doc Document) GetOninvalid() js.Value {
	return js.Value(doc).Get("oninvalid")
}

func (doc Document) SetOninvalid(val any) {
	js.Value(doc).Set("oninvalid", val)
}


func (doc Document) GetOnkeydown() js.Value {
	return js.Value(doc).Get("onkeydown")
}

func (doc Document) SetOnkeydown(val any) {
	js.Value(doc).Set("onkeydown", val)
}


func (doc Document) GetOnkeypress() js.Value {
	return js.Value(doc).Get("onkeypress")
}

func (doc Document) SetOnkeypress(val any) {
	js.Value(doc).Set("onkeypress", val)
}


func (doc Document) GetOnkeyup() js.Value {
	return js.Value(doc).Get("onkeyup")
}

func (doc Document) SetOnkeyup(val any) {
	js.Value(doc).Set("onkeyup", val)
}


func (doc Document) GetOnload() js.Value {
	return js.Value(doc).Get("onload")
}

func (doc Document) SetOnload(val any) {
	js.Value(doc).Set("onload", val)
}


func (doc Document) GetOnloadeddata() js.Value {
	return js.Value(doc).Get("onloadeddata")
}

func (doc Document) SetOnloadeddata(val any) {
	js.Value(doc).Set("onloadeddata", val)
}


func (doc Document) GetOnloadedmetadata() js.Value {
	return js.Value(doc).Get("onloadedmetadata")
}

func (doc Document) SetOnloadedmetadata(val any) {
	js.Value(doc).Set("onloadedmetadata", val)
}


func (doc Document) GetOnloadstart() js.Value {
	return js.Value(doc).Get("onloadstart")
}

func (doc Document) SetOnloadstart(val any) {
	js.Value(doc).Set("onloadstart", val)
}


func (doc Document) GetOnlostpointercapture() js.Value {
	return js.Value(doc).Get("onlostpointercapture")
}

func (doc Document) SetOnlostpointercapture(val any) {
	js.Value(doc).Set("onlostpointercapture", val)
}


func (doc Document) GetOnmousedown() js.Value {
	return js.Value(doc).Get("onmousedown")
}

func (doc Document) SetOnmousedown(val any) {
	js.Value(doc).Set("onmousedown", val)
}


func (doc Document) GetOnmouseenter() js.Value {
	return js.Value(doc).Get("onmouseenter")
}

func (doc Document) SetOnmouseenter(val any) {
	js.Value(doc).Set("onmouseenter", val)
}


func (doc Document) GetOnmouseleave() js.Value {
	return js.Value(doc).Get("onmouseleave")
}

func (doc Document) SetOnmouseleave(val any) {
	js.Value(doc).Set("onmouseleave", val)
}


func (doc Document) GetOnmousemove() js.Value {
	return js.Value(doc).Get("onmousemove")
}

func (doc Document) SetOnmousemove(val any) {
	js.Value(doc).Set("onmousemove", val)
}


func (doc Document) GetOnmouseout() js.Value {
	return js.Value(doc).Get("onmouseout")
}

func (doc Document) SetOnmouseout(val any) {
	js.Value(doc).Set("onmouseout", val)
}


func (doc Document) GetOnmouseover() js.Value {
	return js.Value(doc).Get("onmouseover")
}

func (doc Document) SetOnmouseover(val any) {
	js.Value(doc).Set("onmouseover", val)
}


func (doc Document) GetOnmouseup() js.Value {
	return js.Value(doc).Get("onmouseup")
}

func (doc Document) SetOnmouseup(val any) {
	js.Value(doc).Set("onmouseup", val)
}


func (doc Document) GetOnpaste() js.Value {
	return js.Value(doc).Get("onpaste")
}

func (doc Document) SetOnpaste(val any) {
	js.Value(doc).Set("onpaste", val)
}


func (doc Document) GetOnpause() js.Value {
	return js.Value(doc).Get("onpause")
}

func (doc Document) SetOnpause(val any) {
	js.Value(doc).Set("onpause", val)
}


func (doc Document) GetOnplay() js.Value {
	return js.Value(doc).Get("onplay")
}

func (doc Document) SetOnplay(val any) {
	js.Value(doc).Set("onplay", val)
}


func (doc Document) GetOnplaying() js.Value {
	return js.Value(doc).Get("onplaying")
}

func (doc Document) SetOnplaying(val any) {
	js.Value(doc).Set("onplaying", val)
}


func (doc Document) GetOnpointercancel() js.Value {
	return js.Value(doc).Get("onpointercancel")
}

func (doc Document) SetOnpointercancel(val any) {
	js.Value(doc).Set("onpointercancel", val)
}


func (doc Document) GetOnpointerdown() js.Value {
	return js.Value(doc).Get("onpointerdown")
}

func (doc Document) SetOnpointerdown(val any) {
	js.Value(doc).Set("onpointerdown", val)
}


func (doc Document) GetOnpointerenter() js.Value {
	return js.Value(doc).Get("onpointerenter")
}

func (doc Document) SetOnpointerenter(val any) {
	js.Value(doc).Set("onpointerenter", val)
}


func (doc Document) GetOnpointerleave() js.Value {
	return js.Value(doc).Get("onpointerleave")
}

func (doc Document) SetOnpointerleave(val any) {
	js.Value(doc).Set("onpointerleave", val)
}


func (doc Document) GetOnpointermove() js.Value {
	return js.Value(doc).Get("onpointermove")
}

func (doc Document) SetOnpointermove(val any) {
	js.Value(doc).Set("onpointermove", val)
}


func (doc Document) GetOnpointerout() js.Value {
	return js.Value(doc).Get("onpointerout")
}

func (doc Document) SetOnpointerout(val any) {
	js.Value(doc).Set("onpointerout", val)
}


func (doc Document) GetOnpointerover() js.Value {
	return js.Value(doc).Get("onpointerover")
}

func (doc Document) SetOnpointerover(val any) {
	js.Value(doc).Set("onpointerover", val)
}


func (doc Document) GetOnpointerup() js.Value {
	return js.Value(doc).Get("onpointerup")
}

func (doc Document) SetOnpointerup(val any) {
	js.Value(doc).Set("onpointerup", val)
}


func (doc Document) GetOnprogress() js.Value {
	return js.Value(doc).Get("onprogress")
}

func (doc Document) SetOnprogress(val any) {
	js.Value(doc).Set("onprogress", val)
}


func (doc Document) GetOnratechange() js.Value {
	return js.Value(doc).Get("onratechange")
}

func (doc Document) SetOnratechange(val any) {
	js.Value(doc).Set("onratechange", val)
}


func (doc Document) GetOnreset() js.Value {
	return js.Value(doc).Get("onreset")
}

func (doc Document) SetOnreset(val any) {
	js.Value(doc).Set("onreset", val)
}


func (doc Document) GetOnresize() js.Value {
	return js.Value(doc).Get("onresize")
}

func (doc Document) SetOnresize(val any) {
	js.Value(doc).Set("onresize", val)
}


func (doc Document) GetOnscroll() js.Value {
	return js.Value(doc).Get("onscroll")
}

func (doc Document) SetOnscroll(val any) {
	js.Value(doc).Set("onscroll", val)
}


func (doc Document) GetOnscrollend() js.Value {
	return js.Value(doc).Get("onscrollend")
}

func (doc Document) SetOnscrollend(val any) {
	js.Value(doc).Set("onscrollend", val)
}


func (doc Document) GetOnsecuritypolicyviolation() js.Value {
	return js.Value(doc).Get("onsecuritypolicyviolation")
}

func (doc Document) SetOnsecuritypolicyviolation(val any) {
	js.Value(doc).Set("onsecuritypolicyviolation", val)
}


func (doc Document) GetOnseeked() js.Value {
	return js.Value(doc).Get("onseeked")
}

func (doc Document) SetOnseeked(val any) {
	js.Value(doc).Set("onseeked", val)
}


func (doc Document) GetOnseeking() js.Value {
	return js.Value(doc).Get("onseeking")
}

func (doc Document) SetOnseeking(val any) {
	js.Value(doc).Set("onseeking", val)
}


func (doc Document) GetOnselect() js.Value {
	return js.Value(doc).Get("onselect")
}

func (doc Document) SetOnselect(val any) {
	js.Value(doc).Set("onselect", val)
}


func (doc Document) GetOnselectionchange() js.Value {
	return js.Value(doc).Get("onselectionchange")
}

func (doc Document) SetOnselectionchange(val any) {
	js.Value(doc).Set("onselectionchange", val)
}


func (doc Document) GetOnselectstart() js.Value {
	return js.Value(doc).Get("onselectstart")
}

func (doc Document) SetOnselectstart(val any) {
	js.Value(doc).Set("onselectstart", val)
}


func (doc Document) GetOnslotchange() js.Value {
	return js.Value(doc).Get("onslotchange")
}

func (doc Document) SetOnslotchange(val any) {
	js.Value(doc).Set("onslotchange", val)
}


func (doc Document) GetOnstalled() js.Value {
	return js.Value(doc).Get("onstalled")
}

func (doc Document) SetOnstalled(val any) {
	js.Value(doc).Set("onstalled", val)
}


func (doc Document) GetOnsubmit() js.Value {
	return js.Value(doc).Get("onsubmit")
}

func (doc Document) SetOnsubmit(val any) {
	js.Value(doc).Set("onsubmit", val)
}


func (doc Document) GetOnsuspend() js.Value {
	return js.Value(doc).Get("onsuspend")
}

func (doc Document) SetOnsuspend(val any) {
	js.Value(doc).Set("onsuspend", val)
}


func (doc Document) GetOntimeupdate() js.Value {
	return js.Value(doc).Get("ontimeupdate")
}

func (doc Document) SetOntimeupdate(val any) {
	js.Value(doc).Set("ontimeupdate", val)
}


func (doc Document) GetOntoggle() js.Value {
	return js.Value(doc).Get("ontoggle")
}

func (doc Document) SetOntoggle(val any) {
	js.Value(doc).Set("ontoggle", val)
}


func (doc Document) GetOntouchcancel() js.Value {
	return js.Value(doc).Get("ontouchcancel")
}

func (doc Document) SetOntouchcancel(val any) {
	js.Value(doc).Set("ontouchcancel", val)
}


func (doc Document) GetOntouchend() js.Value {
	return js.Value(doc).Get("ontouchend")
}

func (doc Document) SetOntouchend(val any) {
	js.Value(doc).Set("ontouchend", val)
}


func (doc Document) GetOntouchmove() js.Value {
	return js.Value(doc).Get("ontouchmove")
}

func (doc Document) SetOntouchmove(val any) {
	js.Value(doc).Set("ontouchmove", val)
}


func (doc Document) GetOntouchstart() js.Value {
	return js.Value(doc).Get("ontouchstart")
}

func (doc Document) SetOntouchstart(val any) {
	js.Value(doc).Set("ontouchstart", val)
}


func (doc Document) GetOntransitioncancel() js.Value {
	return js.Value(doc).Get("ontransitioncancel")
}

func (doc Document) SetOntransitioncancel(val any) {
	js.Value(doc).Set("ontransitioncancel", val)
}


func (doc Document) GetOntransitionend() js.Value {
	return js.Value(doc).Get("ontransitionend")
}

func (doc Document) SetOntransitionend(val any) {
	js.Value(doc).Set("ontransitionend", val)
}


func (doc Document) GetOntransitionrun() js.Value {
	return js.Value(doc).Get("ontransitionrun")
}

func (doc Document) SetOntransitionrun(val any) {
	js.Value(doc).Set("ontransitionrun", val)
}


func (doc Document) GetOntransitionstart() js.Value {
	return js.Value(doc).Get("ontransitionstart")
}

func (doc Document) SetOntransitionstart(val any) {
	js.Value(doc).Set("ontransitionstart", val)
}


func (doc Document) GetOnvolumechange() js.Value {
	return js.Value(doc).Get("onvolumechange")
}

func (doc Document) SetOnvolumechange(val any) {
	js.Value(doc).Set("onvolumechange", val)
}


func (doc Document) GetOnwaiting() js.Value {
	return js.Value(doc).Get("onwaiting")
}

func (doc Document) SetOnwaiting(val any) {
	js.Value(doc).Set("onwaiting", val)
}


func (doc Document) GetOnwebkitanimationend() js.Value {
	return js.Value(doc).Get("onwebkitanimationend")
}

func (doc Document) SetOnwebkitanimationend(val any) {
	js.Value(doc).Set("onwebkitanimationend", val)
}


func (doc Document) GetOnwebkitanimationiteration() js.Value {
	return js.Value(doc).Get("onwebkitanimationiteration")
}

func (doc Document) SetOnwebkitanimationiteration(val any) {
	js.Value(doc).Set("onwebkitanimationiteration", val)
}


func (doc Document) GetOnwebkitanimationstart() js.Value {
	return js.Value(doc).Get("onwebkitanimationstart")
}

func (doc Document) SetOnwebkitanimationstart(val any) {
	js.Value(doc).Set("onwebkitanimationstart", val)
}


func (doc Document) GetOnwebkittransitionend() js.Value {
	return js.Value(doc).Get("onwebkittransitionend")
}

func (doc Document) SetOnwebkittransitionend(val any) {
	js.Value(doc).Set("onwebkittransitionend", val)
}


func (doc Document) GetOnwheel() js.Value {
	return js.Value(doc).Get("onwheel")
}

func (doc Document) SetOnwheel(val any) {
	js.Value(doc).Set("onwheel", val)
}


func (doc Document) GetChildElementCount() float64 {
	return js.Value(doc).Get("childElementCount").Float()
}

func (doc Document) SetChildElementCount(val float64) {
	js.Value(doc).Set("childElementCount", val)
}


func (doc Document) GetChildren() js.Value {
	return js.Value(doc).Get("children")
}

func (doc Document) SetChildren(val any) {
	js.Value(doc).Set("children", val)
}


func (doc Document) GetFirstElementChild() js.Value {
	return js.Value(doc).Get("firstElementChild")
}

func (doc Document) SetFirstElementChild(val any) {
	js.Value(doc).Set("firstElementChild", val)
}


func (doc Document) GetLastElementChild() js.Value {
	return js.Value(doc).Get("lastElementChild")
}

func (doc Document) SetLastElementChild(val any) {
	js.Value(doc).Set("lastElementChild", val)
}



func (doc Document) AdoptNode(node any) js.Value {
	return js.Value(doc).Call("adoptNode", node)
}


func (doc Document) CaptureEvents() js.Value {
	return js.Value(doc).Call("captureEvents")
}


func (doc Document) CaretRangeFromPoint(x, y float64) js.Value {
	return js.Value(doc).Call("caretRangeFromPoint", x, y)
}


func (doc Document) Clear() js.Value {
	return js.Value(doc).Call("clear")
}


func (doc Document) Close() js.Value {
	return js.Value(doc).Call("close")
}


func (doc Document) CreateAttribute(localName string) js.Value {
	return js.Value(doc).Call("createAttribute", localName)
}


func (doc Document) CreateAttributeNS(namespace, qualifiedName string) js.Value {
	return js.Value(doc).Call("createAttributeNS", namespace, qualifiedName)
}


func (doc Document) CreateCDATASection(data string) js.Value {
	return js.Value(doc).Call("createCDATASection", data)
}


func (doc Document) CreateComment(data string) js.Value {
	return js.Value(doc).Call("createComment", data)
}


func (doc Document) CreateDocumentFragment() js.Value {
	return js.Value(doc).Call("createDocumentFragment")
}


func (doc Document) CreateElement1K(tagName any) js.Value {
	return js.Value(doc).Call("createElement", tagName)
}


func (doc Document) CreateElement1String(tagName string) js.Value {
	return js.Value(doc).Call("createElement", tagName)
}


func (doc Document) CreateElement2KElementCreationOptions(tagName any, options any) js.Value {
	return js.Value(doc).Call("createElement", tagName, options)
}


func (doc Document) CreateElement2StringElementCreationOptions(tagName string, options any) js.Value {
	return js.Value(doc).Call("createElement", tagName, options)
}


func (doc Document) CreateElementNS2StringString(namespace, qualifiedName string) js.Value {
	return js.Value(doc).Call("createElementNS", namespace, qualifiedName)
}


func (doc Document) CreateElementNS2StringK(namespaceURI string, qualifiedName any) js.Value {
	return js.Value(doc).Call("createElementNS", namespaceURI, qualifiedName)
}


func (doc Document) CreateElementNS2StringString(namespaceURI, qualifiedName string) js.Value {
	return js.Value(doc).Call("createElementNS", namespaceURI, qualifiedName)
}


func (doc Document) CreateElementNS3StringStringString(namespace, qualifiedName, options string) js.Value {
	return js.Value(doc).Call("createElementNS", namespace, qualifiedName, options)
}


func (doc Document) CreateElementNS3StringStringElementCreationOptions(namespaceURI, qualifiedName string, options any) js.Value {
	return js.Value(doc).Call("createElementNS", namespaceURI, qualifiedName, options)
}


func (doc Document) CreateEvent(eventInterface string) js.Value {
	return js.Value(doc).Call("createEvent", eventInterface)
}


func (doc Document) CreateNodeIterator2(root any, whatToShow float64) js.Value {
	return js.Value(doc).Call("createNodeIterator", root, whatToShow)
}


func (doc Document) CreateNodeIterator3(root any, whatToShow float64, filter any) js.Value {
	return js.Value(doc).Call("createNodeIterator", root, whatToShow, filter)
}


func (doc Document) CreateProcessingInstruction(target, data string) js.Value {
	return js.Value(doc).Call("createProcessingInstruction", target, data)
}


func (doc Document) CreateRange() js.Value {
	return js.Value(doc).Call("createRange")
}


func (doc Document) CreateTextNode(data string) js.Value {
	return js.Value(doc).Call("createTextNode", data)
}


func (doc Document) CreateTreeWalker2(root any, whatToShow float64) js.Value {
	return js.Value(doc).Call("createTreeWalker", root, whatToShow)
}


func (doc Document) CreateTreeWalker3(root any, whatToShow float64, filter any) js.Value {
	return js.Value(doc).Call("createTreeWalker", root, whatToShow, filter)
}


func (doc Document) ExecCommand2(commandId string, showUI bool) js.Value {
	return js.Value(doc).Call("execCommand", commandId, showUI)
}


func (doc Document) ExecCommand3(commandId string, showUI bool, value string) js.Value {
	return js.Value(doc).Call("execCommand", commandId, showUI, value)
}


func (doc Document) ExitFullscreen() js.Value {
	return js.Value(doc).Call("exitFullscreen")
}


func (doc Document) ExitPictureInPicture() js.Value {
	return js.Value(doc).Call("exitPictureInPicture")
}


func (doc Document) ExitPointerLock() js.Value {
	return js.Value(doc).Call("exitPointerLock")
}


func (doc Document) GetElementById(elementId string) js.Value {
	return js.Value(doc).Call("getElementById", elementId)
}


func (doc Document) GetElementsByClassName(classNames string) js.Value {
	return js.Value(doc).Call("getElementsByClassName", classNames)
}


func (doc Document) GetElementsByName(elementName string) js.Value {
	return js.Value(doc).Call("getElementsByName", elementName)
}


func (doc Document) GetElementsByTagName1K(qualifiedName any) js.Value {
	return js.Value(doc).Call("getElementsByTagName", qualifiedName)
}


func (doc Document) GetElementsByTagName1String(qualifiedName string) js.Value {
	return js.Value(doc).Call("getElementsByTagName", qualifiedName)
}


func (doc Document) GetElementsByTagNameNS(namespace, localName string) js.Value {
	return js.Value(doc).Call("getElementsByTagNameNS", namespace, localName)
}


func (doc Document) GetSelection() js.Value {
	return js.Value(doc).Call("getSelection")
}


func (doc Document) HasFocus() js.Value {
	return js.Value(doc).Call("hasFocus")
}


func (doc Document) HasStorageAccess() js.Value {
	return js.Value(doc).Call("hasStorageAccess")
}


func (doc Document) ImportNode1(node any) js.Value {
	return js.Value(doc).Call("importNode", node)
}


func (doc Document) ImportNode2(node any, deep bool) js.Value {
	return js.Value(doc).Call("importNode", node, deep)
}


func (doc Document) Open1(unused1 string) js.Value {
	return js.Value(doc).Call("open", unused1)
}


func (doc Document) Open2(unused1, unused2 string) js.Value {
	return js.Value(doc).Call("open", unused1, unused2)
}


func (doc Document) Open3(url, name, features string) js.Value {
	return js.Value(doc).Call("open", url, name, features)
}


func (doc Document) QueryCommandEnabled(commandId string) js.Value {
	return js.Value(doc).Call("queryCommandEnabled", commandId)
}


func (doc Document) QueryCommandIndeterm(commandId string) js.Value {
	return js.Value(doc).Call("queryCommandIndeterm", commandId)
}


func (doc Document) QueryCommandState(commandId string) js.Value {
	return js.Value(doc).Call("queryCommandState", commandId)
}


func (doc Document) QueryCommandSupported(commandId string) js.Value {
	return js.Value(doc).Call("queryCommandSupported", commandId)
}


func (doc Document) QueryCommandValue(commandId string) js.Value {
	return js.Value(doc).Call("queryCommandValue", commandId)
}


func (doc Document) ReleaseEvents() js.Value {
	return js.Value(doc).Call("releaseEvents")
}


func (doc Document) RequestStorageAccess() js.Value {
	return js.Value(doc).Call("requestStorageAccess")
}


func (doc Document) Write(text any) js.Value {
	return js.Value(doc).Call("write", text)
}


func (doc Document) Writeln(text any) js.Value {
	return js.Value(doc).Call("writeln", text)
}


func (doc Document) AddEventListener2K-(this: Document, ev: DocumentEventMap[K]) => any(typeName any, listener any) js.Value {
	return js.Value(doc).Call("addEventListener", typeName, listener)
}


func (doc Document) AddEventListener2StringEventListenerOrEventListenerObject(typeName string, listener any) js.Value {
	return js.Value(doc).Call("addEventListener", typeName, listener)
}


func (doc Document) AddEventListener3K-(this: Document, ev: DocumentEventMap[K]) => anyBoolean(typeName any, listener any, options bool) js.Value {
	return js.Value(doc).Call("addEventListener", typeName, listener, options)
}


func (doc Document) AddEventListener3StringEventListenerOrEventListenerObjectBoolean(typeName string, listener any, options bool) js.Value {
	return js.Value(doc).Call("addEventListener", typeName, listener, options)
}


func (doc Document) RemoveEventListener2K-(this: Document, ev: DocumentEventMap[K]) => any(typeName any, listener any) js.Value {
	return js.Value(doc).Call("removeEventListener", typeName, listener)
}


func (doc Document) RemoveEventListener2StringEventListenerOrEventListenerObject(typeName string, listener any) js.Value {
	return js.Value(doc).Call("removeEventListener", typeName, listener)
}


func (doc Document) RemoveEventListener3K-(this: Document, ev: DocumentEventMap[K]) => anyBoolean(typeName any, listener any, options bool) js.Value {
	return js.Value(doc).Call("removeEventListener", typeName, listener, options)
}


func (doc Document) RemoveEventListener3StringEventListenerOrEventListenerObjectBoolean(typeName string, listener any, options bool) js.Value {
	return js.Value(doc).Call("removeEventListener", typeName, listener, options)
}


func (doc Document) AppendChild(node any) js.Value {
	return js.Value(doc).Call("appendChild", node)
}


func (doc Document) CloneNode0() js.Value {
	return js.Value(doc).Call("cloneNode")
}


func (doc Document) CloneNode1(deep bool) js.Value {
	return js.Value(doc).Call("cloneNode", deep)
}


func (doc Document) CompareDocumentPosition(other any) js.Value {
	return js.Value(doc).Call("compareDocumentPosition", other)
}


func (doc Document) Contains(other any) js.Value {
	return js.Value(doc).Call("contains", other)
}


func (doc Document) GetRootNode0() js.Value {
	return js.Value(doc).Call("getRootNode")
}


func (doc Document) GetRootNode1(options any) js.Value {
	return js.Value(doc).Call("getRootNode", options)
}


func (doc Document) HasChildNodes() js.Value {
	return js.Value(doc).Call("hasChildNodes")
}


func (doc Document) InsertBefore(node any, child any) js.Value {
	return js.Value(doc).Call("insertBefore", node, child)
}


func (doc Document) IsDefaultNamespace(namespace string) js.Value {
	return js.Value(doc).Call("isDefaultNamespace", namespace)
}


func (doc Document) IsEqualNode(otherNode any) js.Value {
	return js.Value(doc).Call("isEqualNode", otherNode)
}


func (doc Document) IsSameNode(otherNode any) js.Value {
	return js.Value(doc).Call("isSameNode", otherNode)
}


func (doc Document) LookupNamespaceURI(prefix string) js.Value {
	return js.Value(doc).Call("lookupNamespaceURI", prefix)
}


func (doc Document) LookupPrefix(namespace string) js.Value {
	return js.Value(doc).Call("lookupPrefix", namespace)
}


func (doc Document) Normalize() js.Value {
	return js.Value(doc).Call("normalize")
}


func (doc Document) RemoveChild(child any) js.Value {
	return js.Value(doc).Call("removeChild", child)
}


func (doc Document) ReplaceChild(node any, child any) js.Value {
	return js.Value(doc).Call("replaceChild", node, child)
}


func (doc Document) DispatchEvent(event any) js.Value {
	return js.Value(doc).Call("dispatchEvent", event)
}


func (doc Document) ElementFromPoint(x, y float64) js.Value {
	return js.Value(doc).Call("elementFromPoint", x, y)
}


func (doc Document) ElementsFromPoint(x, y float64) js.Value {
	return js.Value(doc).Call("elementsFromPoint", x, y)
}


func (doc Document) GetAnimations() js.Value {
	return js.Value(doc).Call("getAnimations")
}


func (doc Document) Append(nodes string) js.Value {
	return js.Value(doc).Call("append", nodes)
}


func (doc Document) Prepend(nodes string) js.Value {
	return js.Value(doc).Call("prepend", nodes)
}


func (doc Document) QuerySelector1K(selectors any) js.Value {
	return js.Value(doc).Call("querySelector", selectors)
}


func (doc Document) QuerySelector1String(selectors string) js.Value {
	return js.Value(doc).Call("querySelector", selectors)
}


func (doc Document) QuerySelectorAll1K(selectors any) js.Value {
	return js.Value(doc).Call("querySelectorAll", selectors)
}


func (doc Document) QuerySelectorAll1String(selectors string) js.Value {
	return js.Value(doc).Call("querySelectorAll", selectors)
}


func (doc Document) ReplaceChildren(nodes string) js.Value {
	return js.Value(doc).Call("replaceChildren", nodes)
}


func (doc Document) CreateExpression1(expression string) js.Value {
	return js.Value(doc).Call("createExpression", expression)
}


func (doc Document) CreateExpression2(expression string, resolver any) js.Value {
	return js.Value(doc).Call("createExpression", expression, resolver)
}


func (doc Document) CreateNSResolver(nodeResolver any) js.Value {
	return js.Value(doc).Call("createNSResolver", nodeResolver)
}


func (doc Document) Evaluate4(expression string, contextNode any, resolver any, typeName float64) js.Value {
	return js.Value(doc).Call("evaluate", expression, contextNode, resolver, typeName)
}


func (doc Document) Evaluate5(expression string, contextNode any, resolver any, typeName float64, result any) js.Value {
	return js.Value(doc).Call("evaluate", expression, contextNode, resolver, typeName, result)
}



type DocumentReadyState = string
const (
  DocumentReadyStateComplete DocumentReadyState = "complete"
  DocumentReadyStateInteractive DocumentReadyState = "interactive"
  DocumentReadyStateLoading DocumentReadyState = "loading"
)


type DocumentVisibilityState = string
const (
  DocumentVisibilityStateHidden DocumentVisibilityState = "hidden"
  DocumentVisibilityStateVisible DocumentVisibilityState = "visible"
)



