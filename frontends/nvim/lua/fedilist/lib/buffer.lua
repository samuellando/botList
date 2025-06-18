local M = {buffers = {}}

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

return M
