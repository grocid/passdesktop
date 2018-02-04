package dialog

func Confirm() string {
  return `
<div class="WindowLayout">
    <div class="SearchLayout">
        <input type="text"
               value="{{html .Title}}"
               placeholder="Account"
               onchange="DoSearchQuery"
               autocomplete="off"
               autocorrect="off"
               autocapitalize="off"
               spellcheck="false"
               selectable="on"
               class="editable searchfield"/>
    </div>
    <div class="animated">
        <div class="symbol trash"/>
        <div class="bottom-toolbar">
            <button class="button ok" onclick="ConfirmOK"/>
            <button class="button cancel" onclick="ConfirmCancel"/>
        </div>
    </div>
</div>`
}
