package main

import (
    "crypto/sha1" // for non-cryptographic purposes
    "encoding/hex"
    "fmt"
    "github.com/murlokswarm/app"
    "io/ioutil"
    "log"
    "net/url"
    "os/user"
    "pass/rest"
)

type File struct {
    Title   string
    Query   string
    Changed bool
    Data    rest.DecodedEntry
}

func (h *File) Render() string {
    // Calculate file size in kBs.
    fs := fmt.Sprintf("%.2f", float64(len(h.Data.File))/1024)
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
        <div class="animated">
            <div style="text-align: center;
                        margin-left: auto;
                        margin-right: auto;
                        margin-top: -webkit-calc(20vh - 20px);">
            ` + GetFingerprint(h.Data.File, 255, 255, 255) + `
            <h1>{{.Title}}</h1>
            </div>
          <h2>Size</h2>
          <p style="text-align: center">` + fs + ` kB</p>
          </div>
          <div class="bottom-toolbar">
              <div>
                  <button class="button ok" onclick="OK"/>
                  <button class="button add" onclick="ReadFile"/>
                  <button class="button download" onclick="SaveFile"/>
                  <button class="button delete" onclick="Delete"/>
              </div>
          </div>
     </div>
</div>`

}

func GetFingerprint(data []byte, r int, g int, b int) string {
    hashFunction := sha1.New()
    hashFunction.Write(data)
    hashDigest := hashFunction.Sum(nil)

    grid := `<div class="grid-container">`
    for i := 0; i < 16; i++ {
        item := fmt.Sprintf(`<div class="grid-item"
                                  style="background-color: rgba(%v, %v, %v, 0.%v)">
                            </div>`, r, g, b, 10*int(hashDigest[i])/255)
        grid = grid + item
    }

    return grid + `</div>`
}

func (h *File) OnHref(URL *url.URL) {
    // Extract information from query and get account name and
    // its encrypted counterpart from query (this is need since
    // if we were to encrypt again, we would get a different
    // encrypted name).
    u := URL.Query()
    h.Title = u.Get("Name")
    restResponse := restClient.VaultReadSecret(
        &rest.Name{
            Text:      h.Title,
            Encrypted: u.Get("Encrypted"),
        })
    h.Data = *restResponse

    // Tells the app to update the rendering of the component.
    app.Render(h)
}

func (h *File) OK() {
    // Make sure we do not save already saved information.
    if h.Changed {
        // No empy names.
        if h.Title == "" {
            return
        }

        d := h.Data.Name
        if h.Title != (*d).Text {
            // need to remove old and submit new
        } else {
            // Modify the decoded entry so that it matches
            //the contents of the UI.
            restClient.VaultWriteSecret(&h.Data)
        }
    }

    // Now, we just need to go back.
    h.Cancel()
}

func (h *File) Cancel() {
    NavigateBack("")
}

func (h *File) SaveFile() {
    usr, err := user.Current()

    if err != nil {
        log.Fatal(err)
    }

    hashFunction := sha1.New()
    hashFunction.Write(h.Data.File)
    hexDigest := hex.EncodeToString(hashFunction.Sum(nil))

    filename := usr.HomeDir + "/Downloads/" + hexDigest
    ioutil.WriteFile(filename, h.Data.File, 0644)

    log.Println("Wrote file" + filename)

    app.Render(h)
}

func (h *File) ReadFile() {
    // Open filepicker window to get filename.
    app.NewFilePicker(app.FilePicker{
        MultipleSelection: false,
        NoDir:             true,
        NoFile:            false,
        OnPick: func(filenames []string) {
            // Get contents of file.
            b, err := ioutil.ReadFile(filenames[0])

            // If there was an error, probably due
            // to permissions, dump error message
            // to log.
            if err != nil {
                log.Println(err)
            } else {
                h.Data.File = b
                h.Changed = true
            }

            app.Render(h)
        },
    })
}

func (h *File) Delete() {
    d := h.Data.Name
    if d != nil {
        restClient.VaultDeleteSecret(&h.Data)
    }
    h.Cancel()
}

func (h *File) DoSearchQuery(arg app.ChangeArg) {
    NavigateBack(arg.Value)
}

func init() {
    app.RegisterComponent(&File{})
}
