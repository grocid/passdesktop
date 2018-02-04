package main

import (
    "github.com/murlokswarm/app"
)

type About struct{}

func (*About) Render() string {
    return `
<div class="WindowLayout">
    <div class="animated">
        <div style="text-align: center;
                    margin: 0 auto;
                    margin-top: -webkit-calc(20vh - 20px);
                    max-width: 360px;">
            <img src="iconpack/default.png" 
                 style="max-width: 128px; "/>
            <h1>Pass Desktop</h1>
            <h2>
                Written by Carl Löndahl<br/>grocid.net
            </h2>
            <p>
                Pass Desktop is free software and licensed under the BSD three-clause (revised) license.
            </p>
            <p>
                Copyright © 2018 Carl Löndahl. All rights reserved.
            </p>
        </div>
        <div class="bottom-toolbar">
            <button class="button ok" onclick="OK"/>
        </div>
    </div>
</div>`
}

func (h *About) OK() {
    if !pass.Locked {
        // If it was unlocked, we can go back to search...
        NavigateBack("")
    } else {
        // otherwise, we need to return to unlock screen.
        s := UnlockScreen{}
        win.Mount(&s)
    }
}

func init() {
    app.RegisterComponent(&About{})
}
