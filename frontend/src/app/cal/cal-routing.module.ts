import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { CalComponent } from './cal.component';


const routes: Routes = [
  {
    path: '',
    component: CalComponent,
    data: {
      title: 'Cal'
    }
    // children: [
    //   {
    //     path: 'page',
    //     component: PageComponent,
    //     data: {
    //       title: 'Page'
    //     }
    //   }
    // ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class CalRoutingModule { }
