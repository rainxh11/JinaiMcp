export class RPCReflect {
    private methods: Map<string, Function> = new Map();

    registerMethod(name: string, fn: Function) {
        this.methods.set(name, fn);
    }

    getMethod(name: string): Function | undefined {
        return this.methods.get(name);
    }

    listMethods(): string[] {
        return Array.from(this.methods.keys());
    }
}
