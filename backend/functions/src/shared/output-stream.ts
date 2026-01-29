export class OutputServerEventStream {
    private writable: any;

    constructor(writable: any) {
        this.writable = writable;
    }

    write(data: any) {
        this.writable.write(data);
    }

    end() {
        this.writable.end();
    }
}
