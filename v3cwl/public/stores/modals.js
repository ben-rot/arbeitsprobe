document.addEventListener("alpine:init", () => {
    // Auto Focus on Open + Focus Next on Enter
    Alpine.store('modals', {

        current: null,

        show: false,
        mouseDownOnBackdrop: false,


        open(name, endpoint) {
            
            if (name === this.current) {
                this.show = true
                return
            }

            this.current = name 
            this.show = true

            htmx.ajax('GET', endpoint, '#modal-slot', { swap: 'innerHTML' })
        },

        close() { this.show = false },


        closeOnClickOutside() {
            if (this.mouseDownOnBackdrop) this.close()
            this.mouseDownOnBackdrop = false
        }
    })




    Alpine.data('modalRoot', () => {

        const m = Alpine.store('modals')

        return {
            directive: {

                ['x-show']() { return m.show },

                ['@keydown.escape.window']() { m.close() },
                ['@mousedown.self']() { m.mouseDownOnBackdrop = true },
                ['@mouseup.self']() { m.closeOnClickOutside() },
                ['@mouseup']() { m.mouseDownOnBackdrop = false },
            },

            closeModal() { m.close() },
        }
    })
})