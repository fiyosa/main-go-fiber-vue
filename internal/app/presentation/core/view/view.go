package view

import _ "embed"

// SPA is the HTML shell served for every frontend route. The Vite build only
// emits index.js, index.css and assets/*, so the shell lives here instead of
// in public/build.
//
//go:embed spa.html
var SPA string
