document.addEventListener("alpine:init", () => {

    const getPath = () => window.location.pathname.split('/')[1] || 'home';
    
    Alpine.store('nav', {
        active: getPath(),
        setActive(name) {
            this.active = name
        }
    })

    document.addEventListener('htmx:historyRestore', (event) => {
        const path = event.detail.path.split('/')[1] || 'home';
        Alpine.store('nav').setActive(path);
    });

    document.addEventListener('htmx:pushedIntoHistory', (event) => {
        const path = event.detail.path.split('/')[1] || 'home';
        Alpine.store('nav').setActive(path);
    });
})