import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {AdblockComponent} from 'src/app/components/adblock/adblock.component';
import {DynamicComponent} from 'src/app/components/dynamic/dynamic.component';
import {FixedComponent} from 'src/app/components/fixed/fixed.component';

const routes: Routes = [
    {
        path: 'cache/dynamic',
        component: DynamicComponent
    },
    {
        path: 'cache/adblock',
        component: AdblockComponent
    },
    {
        path: 'cache/fixed',
        component: FixedComponent
    },
    {
        path: '',
        redirectTo: 'cache/dynamic',
        pathMatch: 'full',
    },
];

@NgModule({
    imports: [RouterModule.forRoot(routes, {useHash: true})],
    exports: [RouterModule]
})
export class AppRoutingModule {
}
