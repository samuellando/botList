local M = {}

function M.setup(opts)
    vim.api.nvim_create_user_command("GetList", function(opts)
        M.display_by_id(opts.args)
    end, { nargs = 1 })
end

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

M.buffers = M.buffers or {}
function M.open_buffer(name, keymaps, on_save)
    if not M.buffers[name] then
        vim.cmd("enew")
        local buf = vim.api.nvim_get_current_buf()
        vim.api.nvim_buf_set_name(buf, name)
        vim.bo[buf].buftype = "acwrite"
        vim.bo[buf].swapfile = false
        -- Set the keymaps for the buffer
        for lhs, rhs in pairs(keymaps) do
            vim.keymap.set('n', lhs, function() rhs(buf) end, { buffer = buf })
        end
        -- Define our custom on save function
        vim.api.nvim_create_autocmd("BufWriteCmd", {
            buffer = buf,
            callback = function()
                on_save(buf)
            end,
        })
        M.buffers[name] = buf
        return buf
    else
        vim.api.nvim_set_current_buf(M.buffers[name])
        return M.buffers[name]
    end
end

function M.display(url)
    local ok, data = pcall(vim.fn.json_decode, fetch_url(url))
    if not ok or not data then
        vim.api.nvim_err_writeln("Invalid JSON")
        return
    end
    -- create/open the buffer
    local name = "fedilist://" .. url
    local buf = M.open_buffer(name, {
        ["<CR>"] = M.on_enter,
        ["g."] = M.on_g_dot,
    }, M.on_save)
    -- Hide the ids at the start of the line
    vim.api.nvim_buf_set_option(buf, "conceallevel", 2)
    vim.api.nvim_buf_set_option(buf, "concealcursor", "n")
    -- Define a syntax region that will be concealed (local to buffer)
    vim.api.nvim_buf_call(buf, function()
        vim.cmd [[syntax match FediListConcealed /^\\\w\{3}\s/ conceal]]
        vim.cmd [[highlight default link FediListConcealed Conceal]]
    end)
    -- Prepare lines
    local lines = {
        data.name,
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
    local matched = {}
    local lines = vim.api.nvim_buf_get_lines(buf, 0, -1, false)
    for _, line in ipairs(lines) do
        local rid = line:match("^\\(%w%w%w)%s")
        local name = line:match("^%- (.+)")
        if not rid then
            if not (rid and ids[rid]) and name then
                M.append(buf, name)
            end
        else
            matched[rid] = true
        end
    end
    for rid, _ in pairs(ids) do
        if not matched[rid] then
            M.remove(buf, ids[rid])
        end
    end
    vim.api.nvim_buf_set_option(buf, "modified", false)
    M.display(vim.b[buf].id)
end

function M.remove(buf, id)
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
        type = "RemoveAction",
        agent = {
            type = "Person",
            id   = "http://localhost:9090/user/samuel"
        },
        object = {
            type = "ItemList",
            id = id
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

function M.append(buf, name)
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

function M.on_enter(buf)
    local line = vim.api.nvim_win_get_cursor(0)[1]
    local txt = vim.api.nvim_buf_get_lines(buf, line - 1, line, false)[1]
    local rid = txt and txt:match("^\\(%w%w%w)%s")
    if rid and vim.b[buf].ids[rid] then
        M.display(vim.b[buf].ids[rid])
    end
end

function M.on_g_dot(parent_buf)
    local url = vim.b[parent_buf].id

    local ok, data = pcall(vim.fn.json_decode, fetch_url(url))
    if not ok or not data then
        vim.api.nvim_err_writeln("Invalid JSON")
        return
    end

    local name = vim.api.nvim_buf_get_name(parent_buf) .. "/raw.json"
    local newbuf = M.open_buffer(name, {}, function() end)
    vim.api.nvim_buf_set_option(newbuf, "filetype", "json")

    data.itemListElement = nil
    data.numberOfItems = nil
    data["@context"] = nil

    local json = vim.fn.json_encode(data)
    local pretty = vim.fn.system({ "jq", "." }, json)

    local lines = {}
    for line in pretty:gmatch("([^\n]+)") do
        table.insert(lines, line)
    end

    vim.api.nvim_buf_set_lines(newbuf, 0, -1, false, lines)
end

return M
