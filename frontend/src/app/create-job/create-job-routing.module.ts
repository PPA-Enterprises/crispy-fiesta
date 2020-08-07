import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { CreateJobComponent } from './create-job.component';


const routes: Routes = [
  {
    path: '',
    component: CreateJobComponent,
    data: {
      title: 'Create Job'
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
export class CreateJobRoutingModule { }
