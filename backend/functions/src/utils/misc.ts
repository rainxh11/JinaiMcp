export function cleanAttribute(attr: any): string {
    if (attr === null || attr === undefined) {
        return '';
    }
    return String(attr).trim();
}
