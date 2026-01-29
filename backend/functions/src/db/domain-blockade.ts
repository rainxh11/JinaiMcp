export class DomainBlockade {
    domain: string;
    triggerReason: string;
    triggerUrl: string;
    createdAt: Date;
    expireAt: Date;

    static from(data: any): DomainBlockade {
        const blockade = new DomainBlockade();
        blockade.domain = data.domain;
        blockade.triggerReason = data.triggerReason;
        blockade.triggerUrl = data.triggerUrl;
        blockade.createdAt = data.createdAt;
        blockade.expireAt = data.expireAt;
        return blockade;
    }

    async save(): Promise<void> {
        // In local deployment, we don't persist to database
        console.log(`Domain blockade saved: ${this.domain}`);
    }
}
