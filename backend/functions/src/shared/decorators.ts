export function CloudHTTPv2() {
    return function (target: any, propertyKey: string, descriptor: PropertyDescriptor) {
        // Decorator implementation
    };
}

export function Ctx() {
    return function (target: any, propertyKey: string, parameterIndex: number) {
        // Decorator implementation
    };
}

export type CloudEvent = {
    data: any;
    context: {
        eventId: string;
        timestamp: string;
        eventType: string;
        resource: string;
    };
};
