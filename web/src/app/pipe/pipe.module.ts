import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {DomainFilterPipe} from 'src/app/pipe/domain-filter.pipe';
import {DomainPipe} from 'src/app/pipe/domain.pipe';
import {IndexPipe} from 'src/app/pipe/index.pipe';


@NgModule({
    declarations: [
        DomainPipe,
        DomainFilterPipe,
        IndexPipe
    ],
    exports: [
        DomainPipe,
        DomainFilterPipe,
        IndexPipe
    ],
    imports: [
        CommonModule
    ]
})
export class PipeModule {
}
