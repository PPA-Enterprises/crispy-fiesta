import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { JobsComponent } from './jobs.component';


const routes: Routes = [
  {
    path: '',
    component: JobsComponent,
    data: {
      title: 'Jobs'
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
export class JobsRoutingModule { }
