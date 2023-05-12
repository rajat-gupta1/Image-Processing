#!/bin/bash
#
#SBATCH --mail-user=rajatgupta@uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=part6
#SBATCH --output=/home/rajatgupta/course/parallel/project-2-rajat-gupta1/proj2/out/%j.%N.stdout
#SBATCH --error=/home/rajatgupta/course/parallel/project-2-rajat-gupta1/proj2/out/%j.%N.stderr
#SBATCH --chdir=/home/rajatgupta/course/parallel/project-2-rajat-gupta1/proj2
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=03:00:00


module load golang/1.16.2

# Sequential Small
for i in {1..5}
do
    go run editor/editor.go small >> time.txt
done

# Pipeline Small
for j in 1 2 3 4 6
do
    for i in {1..5}
    do
        go run editor/editor.go small pipeline $(($j * 2)) >> time.txt
    done
done

# Sequential Mixture
for i in {1..5}
do
    go run editor/editor.go mixture >> time.txt
done

# Pipeline Mixture
for j in 1 2 3 4 6
do
    for i in {1..5}
    do
        go run editor/editor.go mixture pipeline $(($j * 2)) >> time.txt
    done
done

# Sequential Big
for i in {1..5}
do
    go run editor/editor.go big >> time.txt
done

# Pipeline Big
for j in 1 2 3 4 6
do
    for i in {1..5}
    do
        go run editor/editor.go big pipeline $(($j * 2)) >> time.txt
    done
done


# BSP Small
for j in 1 2 3 4 6
do
    for i in {1..5}
    do
        go run editor/editor.go small bsp $(($j * 2)) >> time.txt
    done
done

# BSP Mixture
for j in 1 2 3 4 6
do
    for i in {1..5}
    do
        go run editor/editor.go mixture bsp $(($j * 2)) >> time.txt
    done
done

# BSP Big
for j in 1 2 3 4 6
do
    for i in {1..5}
    do
        go run editor/editor.go big bsp $(($j * 2)) >> time.txt
    done
done