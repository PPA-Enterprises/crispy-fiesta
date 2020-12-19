import { Routes, RouterModule } from '@angular/router';

//Route for content layout with sidebar, navbar and footer.

export const Full_ROUTES: Routes = [
  {
    path: 'home',
    loadChildren: () => import('../../home/home.module').then(m => m.HomeModule)
  },
  {
    path: 'jobs',
    loadChildren: () => import('../../jobs/jobs.module').then(m => m.JobsModule)
  },
  // {
  //   path: 'clients',
  //   loadChildren: () => import('../../clients/clients.module').then(m => m.ClientsModule)
  // },
  {
    path: 'create-job',
    loadChildren: () => import('../../create-job/create-job.module').then(m => m.CreateJobModule)
  },
  {
    path: 'edit-job',
    loadChildren: () => import('../../edit-job/edit-job.module').then(m => m.EditJobModule)
  },
];
