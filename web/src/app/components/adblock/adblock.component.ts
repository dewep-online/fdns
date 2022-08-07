import {Component, NgZone, OnInit} from '@angular/core';
import {RequestService} from '@deweppro/core';
import {AdblockDomain, AdblockURI} from 'src/app/models/adblock';
import {CacheItem} from 'src/app/models/cache';
import {environment} from 'src/environments/environment';

@Component({
    selector: 'app-adblock',
    templateUrl: './adblock.component.html',
    styleUrls: ['./adblock.component.scss']
})
export class AdblockComponent implements OnInit {

    uris: AdblockURI[] = [];
    domains: AdblockDomain[] = [];
    filter = '';
    constructor(
        protected http: RequestService
    ) {
    }

    ngOnInit(): void {
        this.reload();
    }

    reload(): void {
        this.uriList();
        this.domainList();
    }

    uriList(): void {
        this.http.get(environment.adblock_list_uri)
            .subscribe((value: any) => {
                this.uris = value;
            });
    }

    domainList(): void {
        this.http.get(environment.adblock_list_domain)
            .subscribe((value: any) => {
                this.domains = value;
            });
    }

    changeActive(domain: string, active: boolean): void {
        active = !active;
        this.http.post(environment.adblock_active, {domain, active})
            .subscribe(() => {
                this.domains = this.domains.map((value) => {
                    if (value.domain === domain) {
                        value.active = active;
                    }
                    return value;
                });
            });
    }

}
