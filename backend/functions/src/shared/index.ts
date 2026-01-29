import { injectable } from 'tsyringe';
import * as fs from 'fs';
import * as path from 'path';

@injectable()
export class AsyncContext {
    private storage: Map<string, any> = new Map();
    set(key: string, value: any) {
        this.storage.set(key, value);
    }

    get(key: string): any {
        return this.storage.get(key);
    }
}

export class InsufficientBalanceError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'InsufficientBalanceError';
    }
}

@injectable()
export class FirebaseStorageBucketControl {
    private localStorageDir: string;

    constructor() {
        this.localStorageDir = path.join('/app', 'local-storage');
        if (!fs.existsSync(this.localStorageDir)) {
            fs.mkdirSync(this.localStorageDir, { recursive: true });
        }
    }

    async uploadFile(filePath: string, destination: string): Promise<string> {
        const destPath = path.join(this.localStorageDir, destination);
        await fs.promises.copyFile(filePath, destPath);
        return `file://${destPath}`;
    }

    async downloadFile(filePath: string, destination: string): Promise<void> {
        const sourcePath = path.join(this.localStorageDir, filePath);
        await fs.promises.copyFile(sourcePath, destination);
    }

    async deleteFile(filePath: string): Promise<void> {
        const fullPath = path.join(this.localStorageDir, filePath);
        await fs.promises.unlink(fullPath);
    }

    async fileExists(filePath: string): Promise<boolean> {
        const fullPath = path.join(this.localStorageDir, filePath);
        return fs.existsSync(fullPath);
    }

    async saveFile(filePath: string, content: Buffer, options?: any): Promise<void> {
        const fullPath = path.join(this.localStorageDir, filePath);
        await fs.promises.writeFile(fullPath, content);
    }

    async signDownloadUrl(filePath: string, expirationTime: number): Promise<string> {
        const fullPath = path.join(this.localStorageDir, filePath);
        return `file://${fullPath}`;
    }
}

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
