import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ErrorComponent } from './pages/error/error.component';
import { HomeComponent } from './pages/home/home.component';

const routes: Routes = [
  { path:'home', component: HomeComponent },
  { path: '', redirectTo: '/home', pathMatch: 'full' },
  { path:'**', component: ErrorComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule],
})
export class AppRoutingModule { }
