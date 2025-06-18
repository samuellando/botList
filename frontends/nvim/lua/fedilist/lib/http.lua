local M = {}

function M.fetch_url(url)
    local handle = io.popen("curl -sL " .. vim.fn.shellescape(url))
    if not handle then return nil end
    local result = handle:read("*a")
    handle:close()
    return result
end

return M
