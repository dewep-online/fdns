import {Component, OnInit} from '@angular/core';
import {RequestService} from '@deweppro/core';
import {CacheItem} from 'src/app/models/cache';
import {environment} from 'src/environments/environment';

@Component({
    selector: 'app-dynamic',
    templateUrl: './dynamic.component.html',
    styleUrls: ['./dynamic.component.scss']
})
export class DynamicComponent implements OnInit {

    list: CacheItem[] = [];
    filter = '';

    constructor(
        protected http: RequestService
    ) {
    }

    ngOnInit(): void {
        this.cacheList();
    }

    cacheList(): void {
        this.http.get(environment.cache_list, {type: 1, filter: this.filter})
            .subscribe((value: any) => {
                this.list = value;
            });
    }

    blockDomain(domain: string): void {
        this.http.post(environment.cache_block, {domain})
            .subscribe(() => {
                this.list = this.list.filter(value => {
                    return value.domain !== domain;
                });
            });
    }
}
