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

function iso_datetime()
    return os.date("!%Y-%m-%dT%H:%M:%SZ")
end

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

function random_id(len)
    local res = {}
    for i = 1, len do
        local v = math.random(0, 35)
        res[i] = string.char((v < 10) and (48 + v) or (55 + v))
    end
    return table.concat(res)
end

M.url_to_bufnr = M.url_to_bufnr or {}
function M.display(url)
    local ok, data = pcall(vim.fn.json_decode, fetch_url(url))
    if not ok or not data then
        vim.api.nvim_err_writeln("Invalid JSON")
        return
    end

    -- Open a new scratch buffer
    -- Set buffer as unlisted and scratch
    if not M.url_to_bufnr[url] then
        vim.cmd("enew")
        local buf = vim.api.nvim_get_current_buf()
        vim.bo[buf].buftype = "acwrite"
        vim.bo[buf].swapfile = false
        vim.api.nvim_buf_set_name(buf, "fedilist://" .. url)
        M.url_to_bufnr[url] = buf
        vim.api.nvim_buf_set_keymap(buf, "n", "<CR>", "<cmd>lua require('fedilist').on_enter(" .. buf .. ")<CR>",
            { noremap = true, silent = true })
        -- Hide the ids at the start of the line
        vim.api.nvim_buf_set_option(buf, "conceallevel", 2)
        vim.api.nvim_buf_set_option(buf, "concealcursor", "n")
        -- Define a syntax region that will be concealed (local to buffer)
        vim.api.nvim_buf_call(buf, function()
            vim.cmd [[syntax match FediListConcealed /^\\\w\{3}\s/ conceal]]
            vim.cmd [[highlight default link FediListConcealed Conceal]]
        end)
        -- Define our custom on save function
        vim.api.nvim_create_autocmd("BufWriteCmd", {
            buffer = buf,
            callback = function()
                M.on_save(buf)
            end,
        })
    else
        vim.api.nvim_set_current_buf(M.url_to_bufnr[url])
    end
    local buf = vim.api.nvim_get_current_buf()

    -- Prepare lines
    local lines = {
        "Name: " .. data.name,
        "--------------------------------"
    }
    local ids = {}
    if data.itemListElement then
        for _, item in ipairs(data.itemListElement) do
            local ref = random_id(3)
            -- Avoid collisions
            while ids[ref] do
                ref = random_id(3)
            end
            ids[ref] = item.id
            table.insert(lines, "\\" .. ref .. " - " .. (item.name or "<no name>"))
        end
    end
    vim.b[buf] = vim.b[buf] or {}
    vim.b[buf].ids = ids
    vim.b[buf].id = url
    -- Write lines to buffer
    vim.api.nvim_buf_set_lines(buf, 0, -1, false, lines)
    vim.api.nvim_buf_set_option(buf, "modified", false)
end

function M.on_save(buf)
    local ids = vim.b[buf].ids or {}
    local lines = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
    for _, line in ipairs(lines) do
        local rid = line:match("^\\(%w%w%w)%s")
        local name = line:match("^%- (.+)")
        if not (rid and ids[rid]) and name then
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
                type = "AppendAction",
                agent = {
                    type = "Person",
                    id   = "http://localhost:9090/user/samuel"
                },
                object = {
                    type = "ItemList",
                    name = name
                },
                targetCollection = {
                    type = "ItemList",
                    id   = vim.b[buf].id
                },
                startTime = iso_datetime()
            }
            local body = vim.fn.json_encode(payload)
            local cmd = 'curl -s -X POST -H "Content-Type: application/json" --data ' ..
                vim.fn.shellescape(body) .. ' http://localhost:9090/user/samuel/outbox'
            os.execute(cmd)
        end
    end
    vim.api.nvim_buf_set_option(buf, "modified", false)
    M.display(vim.b[buf].id)
end

function M.on_enter(buf)
    local line = vim.api.nvim_win_get_cursor(0)[1]
    local txt = vim.api.nvim_buf_get_lines(buf, line - 1, line, false)[1]
    local rid = txt and txt:match("^\\(%w%w%w)%s")
    if rid and vim.b[buf].ids[rid] then
        M.display(vim.b[buf].ids[rid])
    end
end

return M
