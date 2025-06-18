local M = {}

local http = require("fedilist.lib.http")
local buffer = require("fedilist.lib.buffer")
local date = require("fedilist.lib.date")

function M.update(buf)
    local lines = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
    local json_str = table.concat(lines, "\n")
    local updated_list = vim.fn.json_decode(json_str)
    local payload = {
        ["@context"] = {
            "https://schema.org",
            {
                owner   = "http://fedilist.com/owner",
                editor  = "http://fedilist.com/editor",
                viewer  = "http://fedilist.com/viewer",
                atIndex = "http://fedilist.com/toIndex",
                Result  = "http://fedilist.com/Result"
            }
        },
        type = "UpdateAction",
        agent = {
            type = "Person",
            id   = "http://localhost:9090/user/samuel"
        },
        object = updated_list,
        targetCollection = {
            type = "ItemList",
            id   = updated_list.id
        },
        startTime = date.iso_datetime()
    }
    local body = vim.fn.json_encode(payload)
    local cmd = 'curl -s -X POST -H "Content-Type: application/json" --data ' ..
        vim.fn.shellescape(body) .. ' http://localhost:9090/user/samuel/outbox'
    os.execute(cmd)
    vim.api.nvim_buf_set_option(buf, "modified", false)
end

function M.display(parent_buf)
    local url = vim.b[parent_buf].id

    local ok, data = pcall(vim.fn.json_decode, http.fetch_url(url))
    if not ok or not data then
        vim.api.nvim_err_writeln("Invalid JSON")
        return
    end

    local name = vim.api.nvim_buf_get_name(parent_buf) .. "/raw.json"
    local newbuf = buffer.open_buffer(name, {}, M.update)
    vim.api.nvim_buf_set_option(newbuf, "filetype", "json")

    data["@context"] = nil

    local json = vim.fn.json_encode(data)
    local pretty = vim.fn.system({ "jq", "." }, json)

    local lines = {}
    for line in pretty:gmatch("([^\n]+)") do
        table.insert(lines, line)
    end

    vim.api.nvim_buf_set_lines(newbuf, 0, -1, false, lines)
    vim.api.nvim_buf_set_option(newbuf, "modified", false)
end

return M
