export namespace config {
	
	export class AppConfig {
	    interface_name: string;
	    promiscuous: boolean;
	    snap_len: number;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.interface_name = source["interface_name"];
	        this.promiscuous = source["promiscuous"];
	        this.snap_len = source["snap_len"];
	    }
	}

}

export namespace main {
	
	export class HealthStatus {
	    is_root: boolean;
	    pcap_version: string;
	    can_sniff: boolean;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new HealthStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_root = source["is_root"];
	        this.pcap_version = source["pcap_version"];
	        this.can_sniff = source["can_sniff"];
	        this.error = source["error"];
	    }
	}
	export class Stats {
	    total_packets: number;
	    total_alerts: number;
	    protocol_dist: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_packets = source["total_packets"];
	        this.total_alerts = source["total_alerts"];
	        this.protocol_dist = source["protocol_dist"];
	    }
	}

}

export namespace network {
	
	export class NetworkInterface {
	    name: string;
	    description: string;
	    ip_addresses: string[];
	
	    static createFrom(source: any = {}) {
	        return new NetworkInterface(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.ip_addresses = source["ip_addresses"];
	    }
	}

}

