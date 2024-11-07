type ClickOutsideEvent = MouseEvent | TouchEvent;

export function clickOutside(node: HTMLElement, onClickOutside: () => void) {
    const handleClick = (event: ClickOutsideEvent) => {
        if (
            node &&
            !node.contains(event.target as Node) &&
            !event.defaultPrevented
        ) {
            onClickOutside();
        }
    };

    document.addEventListener('click', handleClick, true);
    document.addEventListener('touchstart', handleClick, true);

    return {
        destroy() {
            document.removeEventListener('click', handleClick, true);
            document.removeEventListener('touchstart', handleClick, true);
        }
    };
}
