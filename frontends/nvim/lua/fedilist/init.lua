local M = {}

function M.setup(opts)
    vim.api.nvim_create_user_command("MyHello", function()
        print("Hello from my plugin!")
    end, {})
    vim.api.nvim_create_user_command("GetList", function(opts)
        M.display_by_id(opts.args)
    end, { nargs = 1 })
end

local ns_id = vim.api.nvim_create_namespace('fedilist')

-- Fetch JSON over HTTP using curl
function fetch_url(url)
    local handle = io.popen("curl -sL " .. vim.fn.shellescape(url))
    if not handle then return nil end
    local result = handle:read("*a")
    handle:close()
    return result
end

function M.display_by_id(id)
    M.display("http://localhost:9090/list/" .. id)
end

function M.display(url)
    local ok, data = pcall(vim.fn.json_decode, fetch_url(url))
    if not ok or not data then
        vim.api.nvim_err_writeln("Invalid JSON")
        return
    end

    -- Open a new scratch buffer
    vim.cmd("enew")
    local buf = vim.api.nvim_get_current_buf()

    -- Set buffer as unlisted and scratch
    vim.bo[buf].buftype = "nofile"
    vim.bo[buf].swapfile = false

    vim.api.nvim_buf_set_keymap(buf, "n", "<CR>", "<cmd>lua require('fedilist').on_enter(" .. buf .. ")<CR>",
        { noremap = true, silent = true })

    -- Prepare lines
    local lines = {
        "Name: " .. data.name,
        "--------------------------------"
    }
    if data.itemListElement then
        for _, item in ipairs(data.itemListElement) do
            table.insert(lines, "- " .. (item.name or "<no name>"))
        end
    end

    -- Write lines to buffer
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
    local ids = {}
    if data.itemListElement then
        for i, item in ipairs(data.itemListElement) do
            local eid = vim.api.nvim_buf_set_extmark(buf, ns_id, i + 2, 0, {
                right_gravity = false
            })
            ids[eid] = item.id
        end
    end
    vim.b[buf] = vim.b[buf] or {}
    vim.b[buf].ids = ids
end

function M.on_enter(buf)
    local line = vim.api.nvim_win_get_cursor(0)[1]
    local m = vim.api.nvim_buf_get_extmarks(buf, ns_id, { line, 0 }, { line, -1 }, { details = true })[1]
    if m then
        url = vim.b[buf].ids[m[1]]
        M.display(url)
    end
end

return M
