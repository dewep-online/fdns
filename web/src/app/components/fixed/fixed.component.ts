import {Component, OnInit} from '@angular/core';
import {RequestService} from '@deweppro/core';
import {ToastrService} from 'ngx-toastr';
import {Fixed} from 'src/app/models/fixed';
import {environment} from 'src/environments/environment';

@Component({
    selector: 'app-fixed',
    templateUrl: './fixed.component.html',
    styleUrls: ['./fixed.component.scss']
})
export class FixedComponent implements OnInit {

    list: Fixed[] = [];
    id = 0;

    constructor(
        private http: RequestService,
        private alert: ToastrService
    ) {
    }

    ngOnInit(): void {
        this.rulesList();
    }

    rulesList(): void {
        this.http.get(environment.fixed_list)
            .subscribe((value: any) => {
                this.list = value;
            });
    }

    show(i: number): void {
        this.id = i;
    }

    save(rule: Fixed): void {
        this.http.post(environment.fixed_save, rule)
            .subscribe((value: any) => {
                this.list.splice(this.id, 1, value);
                this.rebuild();
                this.alert.success('Saved!');
            });
    }

    add(): void {
        const el: Fixed = {active: false, origin: '', domain: '', types: '', ips: ''};
        this.list = [el, ...this.list];
        this.show(0);
    }

    del(rule: Fixed): void {
        if (!confirm('Do you really want to delete?')) {
            return;
        }
        if (rule.origin.length === 0) {
            this.list.splice(this.id, 1);
            this.rebuild();
            this.alert.success('Deleted!');
            return;
        }
        this.http.post(environment.fixed_delete, rule)
            .subscribe(() => {
                this.list.splice(this.id, 1);
                this.rebuild();
                this.alert.success('Deleted!');
            });
    }

    active(rule: Fixed): void {
        rule.active = !rule.active;
        this.http.post(environment.fixed_active, rule)
            .subscribe((value: any) => {
                this.list.splice(this.id, 1, value);
                this.rebuild();
                this.alert.success(rule.active ? 'Activated!' : 'Deactivated!');
            });
    }

    private rebuild(): void {
        this.list = [...this.list];
    }
}
