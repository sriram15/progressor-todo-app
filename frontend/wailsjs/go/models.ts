export namespace database {
	
	export class GetCardRow {
	    card_id: number;
	    title: string;
	    description: sql.NullString;
	    createdat: sql.NullTime;
	    updatedat: sql.NullTime;
	    status: number;
	    completedat: sql.NullTime;
	    isactive: boolean;
	    estimatedmins: number;
	    trackedmins: number;
	    time_entry_id: sql.NullInt64;
	    starttime: sql.NullTime;
	    endtime: sql.NullTime;
	
	    static createFrom(source: any = {}) {
	        return new GetCardRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.card_id = source["card_id"];
	        this.title = source["title"];
	        this.description = this.convertValues(source["description"], sql.NullString);
	        this.createdat = this.convertValues(source["createdat"], sql.NullTime);
	        this.updatedat = this.convertValues(source["updatedat"], sql.NullTime);
	        this.status = source["status"];
	        this.completedat = this.convertValues(source["completedat"], sql.NullTime);
	        this.isactive = source["isactive"];
	        this.estimatedmins = source["estimatedmins"];
	        this.trackedmins = source["trackedmins"];
	        this.time_entry_id = this.convertValues(source["time_entry_id"], sql.NullInt64);
	        this.starttime = this.convertValues(source["starttime"], sql.NullTime);
	        this.endtime = this.convertValues(source["endtime"], sql.NullTime);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListCardsRow {
	    id: number;
	    title: string;
	    description: sql.NullString;
	    createdat: sql.NullTime;
	    updatedat: sql.NullTime;
	    status: number;
	    completedat: sql.NullTime;
	    estimatedmins: number;
	    trackedmins: number;
	    isactive: boolean;
	    projectid: number;
	    card_id: number;
	
	    static createFrom(source: any = {}) {
	        return new ListCardsRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = this.convertValues(source["description"], sql.NullString);
	        this.createdat = this.convertValues(source["createdat"], sql.NullTime);
	        this.updatedat = this.convertValues(source["updatedat"], sql.NullTime);
	        this.status = source["status"];
	        this.completedat = this.convertValues(source["completedat"], sql.NullTime);
	        this.estimatedmins = source["estimatedmins"];
	        this.trackedmins = source["trackedmins"];
	        this.isactive = source["isactive"];
	        this.projectid = source["projectid"];
	        this.card_id = source["card_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ListOpenOrCTCardsRow {
	    id: number;
	    title: string;
	    description: sql.NullString;
	    createdat: sql.NullTime;
	    updatedat: sql.NullTime;
	    status: number;
	    completedat: sql.NullTime;
	    estimatedmins: number;
	    trackedmins: number;
	    isactive: boolean;
	    projectid: number;
	    card_id: number;
	
	    static createFrom(source: any = {}) {
	        return new ListOpenOrCTCardsRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = this.convertValues(source["description"], sql.NullString);
	        this.createdat = this.convertValues(source["createdat"], sql.NullTime);
	        this.updatedat = this.convertValues(source["updatedat"], sql.NullTime);
	        this.status = source["status"];
	        this.completedat = this.convertValues(source["completedat"], sql.NullTime);
	        this.estimatedmins = source["estimatedmins"];
	        this.trackedmins = source["trackedmins"];
	        this.isactive = source["isactive"];
	        this.projectid = source["projectid"];
	        this.card_id = source["card_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace service {
	
	export class UpdateCardParams {
	    title: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCardParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.description = source["description"];
	    }
	}

}

export namespace sql {
	
	export class NullInt64 {
	    Int64: number;
	    Valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NullInt64(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Int64 = source["Int64"];
	        this.Valid = source["Valid"];
	    }
	}
	export class NullString {
	    String: string;
	    Valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NullString(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.String = source["String"];
	        this.Valid = source["Valid"];
	    }
	}
	export class NullTime {
	    // Go type: time
	    Time: any;
	    Valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NullTime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Time = this.convertValues(source["Time"], null);
	        this.Valid = source["Valid"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

