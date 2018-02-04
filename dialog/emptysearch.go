package dialog

func EmptySearch() string {
  return `
<div class="WindowLayout">
    <div class="SearchLayout">
        <input type="text"
               value="{{html .Query}}"
               placeholder="Account"
               onchange="Search"
               autocomplete="off"
               autocorrect="off"
               autocapitalize="off"
               spellcheck="false"
               selectable="on"
               class="editable searchfield"/>
    </div>
    <div class="animated">
        <div class="symbol search"/>
        <div class="bottom-toolbar">
            <button class="button add" onclick="AccountCreate"/>
        </div>
    </div>
</div>`
}
