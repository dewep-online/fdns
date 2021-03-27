import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { CacheComponent } from './components/cache/cache.component';

const routes: Routes = [
  {
    path: "cache",
    component: CacheComponent,
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
