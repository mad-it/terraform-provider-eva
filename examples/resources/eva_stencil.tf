resource "eva_stencil" "test" {
    name            = "my stencil"
    organization_id = 1
    language_id     = 1
    country_id      = 1
    header          = "<header></header>"
    template        = "<html></html>"
    footer          = "<footer></footer>"
    helpers         = "function someJavascriptFunction() {}; someJavascriptFunction();"
    type            = 1
    layout          = 1
    destination     = 1
    paper_properties = {
        wait_for_network_idle         = true
        wait_for_js                   = true
        format                        = 1
        orientation                   = 1
        thermal_printer_template_type = 1
        size = {
            width               = "100"
            height              = "500"
            device_scale_factor = 15
        }
        margin = {
            top    = 1
            left   = 1
            bottom = 1
            right  = 1
        }
    }
}