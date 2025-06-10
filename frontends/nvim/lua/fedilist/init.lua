local M = {}

function M.setup(opts)
    vim.api.nvim_create_user_command("MyHello", function()
        print("Hello from my plugin!")
    end, {})
    vim.api.nvim_create_user_command("GetList", function(opts)
        M.display(opts.args)
    end, { nargs = 1 })
end

-- Fetch JSON over HTTP using curl
function fetch_url(url)
    local handle = io.popen("curl -sL " .. vim.fn.shellescape(url))
    if not handle then return nil end
    local result = handle:read("*a")
    handle:close()
    return result
end

function M.get_list(id)
    return fetch_url("http://localhost:9090/list/" .. id)
end

function M.display(id)
    local ok, data = pcall(vim.fn.json_decode, M.get_list(id))
    if not ok or not data then
        vim.api.nvim_err_writeln("Invalid JSON")
        return
    end

    -- Open a new scratch buffer
    vim.cmd("enew")
    local buf = vim.api.nvim_get_current_buf()

    -- Set buffer as unlisted and scratch
    vim.bo[buf].buftype = "nofile"
    vim.bo[buf].bufhidden = "wipe"
    vim.bo[buf].swapfile = false

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
end

return M
