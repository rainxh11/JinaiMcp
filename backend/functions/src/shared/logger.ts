export class Logger {
    constructor(private name: string) {}

    info(message: string, meta?: any) {
        console.log(`[${this.name}] INFO: ${message}`, meta ? JSON.stringify(meta) : '');
    }

    warn(message: string, meta?: any) {
        console.warn(`[${this.name}] WARN: ${message}`, meta ? JSON.stringify(meta) : '');
    }

    error(message: string, meta?: any) {
        console.error(`[${this.name}] ERROR: ${message}`, meta ? JSON.stringify(meta) : '');
    }
}
