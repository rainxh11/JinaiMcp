export interface Ctx {
    req: any;
    res: any;
}

export interface RPCReflection {
    methodName: string;
    parameters: any[];
}
